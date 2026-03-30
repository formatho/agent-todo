package services

import (
	"errors"
	"time"

	"github.com/formatho/agent-todo/backend/config"
	"github.com/formatho/agent-todo/backend/models"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/subscription"
	"gorm.io/gorm"
)

// SubscriptionTier represents different pricing tiers
type SubscriptionTier string

const (
	TierFree      SubscriptionTier = "free"
	TierStarter   SubscriptionTier = "starter"    // $19/mo
	TierPro       SubscriptionTier = "pro"        // $49/mo
	TierEnterprise SubscriptionTier = "enterprise" // $199/mo
)

// TierLimits defines limits for each tier
type TierLimits struct {
	AgentsLimit  int
	OrgsLimit    int
	Features     []string
	Priority     string
	SupportLevel string
}

// GetTierLimits returns limits for a specific tier
func GetTierLimits(tier SubscriptionTier) TierLimits {
	limits := map[SubscriptionTier]TierLimits{
		TierFree: {
			AgentsLimit:  3,
			OrgsLimit:    1,
			Features:     []string{"basic_tasks", "basic_agents"},
			Priority:     "normal",
			SupportLevel: "community",
		},
		TierStarter: {
			AgentsLimit:  10,
			OrgsLimit:    3,
			Features:     []string{"advanced_tasks", "basic_analytics", "email_support"},
			Priority:     "normal",
			SupportLevel: "email",
		},
		TierPro: {
			AgentsLimit:  -1, // Unlimited
			OrgsLimit:    -1, // Unlimited
			Features:     []string{"unlimited_agents", "advanced_analytics", "priority_support", "custom_integrations", "api_access"},
			Priority:     "high",
			SupportLevel: "priority",
		},
		TierEnterprise: {
			AgentsLimit:  -1,
			OrgsLimit:    -1,
			Features:     []string{"everything_in_pro", "sso_saml", "advanced_security", "dedicated_support", "white_label", "sla"},
			Priority:     "highest",
			SupportLevel: "dedicated",
		},
	}
	return limits[tier]
}

// BillingService handles all billing and subscription operations
type BillingService struct {
	db     *gorm.DB
	config *config.Config
}

// NewBillingService creates a new billing service
func NewBillingService(db *gorm.DB, cfg *config.Config) *BillingService {
	// Initialize Stripe
	stripe.Key = cfg.StripeSecretKey

	return &BillingService{
		db:     db,
		config: cfg,
	}
}

// CreateUserSubscription creates a new subscription for a user
func (s *BillingService) CreateUserSubscription(userID, email string, tier SubscriptionTier) (*models.UserSubscription, error) {
	// Check if user already has a subscription
	var existingSub models.UserSubscription
	if err := s.db.Where("user_id = ? AND status = ?", userID, "active").First(&existingSub).Error; err == nil {
		return nil, errors.New("user already has an active subscription")
	}

	// For free tier, create subscription without Stripe
	if tier == TierFree {
		sub := &models.UserSubscription{
			ID:              uuid.New().String(),
			UserID:          userID,
			Tier:            string(tier),
			Status:          "active",
			AgentsLimit:     3,
			OrganizationsLimit: 1,
		}

		if err := s.db.Create(sub).Error; err != nil {
			return nil, err
		}

		return sub, nil
	}

	// For paid tiers, create Stripe checkout session
	checkoutSession, err := s.createCheckoutSession(email, tier)
	if err != nil {
		return nil, err
	}

	// Create pending subscription
	sub := &models.UserSubscription{
		ID:              uuid.New().String(),
		UserID:          userID,
		Tier:            string(tier),
		Status:          "pending",
		StripeSessionID: checkoutSession.ID,
		AgentsLimit:     GetTierLimits(tier).AgentsLimit,
		OrganizationsLimit: GetTierLimits(tier).OrgsLimit,
	}

	if err := s.db.Create(sub).Error; err != nil {
		return nil, err
	}

	return sub, nil
}

// createCheckoutSession creates a Stripe checkout session
func (s *BillingService) createCheckoutSession(email string, tier SubscriptionTier) (*stripe.CheckoutSession, error) {
	// Get price ID based on tier
	priceID := s.getPriceID(tier)

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String("subscription"),
		CustomerEmail:      stripe.String(email),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(s.config.FrontendURL + "/checkout/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(s.config.FrontendURL + "/pricing"),
		Metadata: map[string]string{
			"tier": string(tier),
		},
	}

	return session.New(params)
}

// getPriceID returns the Stripe price ID for a tier
func (s *BillingService) getPriceID(tier SubscriptionTier) string {
	// In production, these would come from config
	priceIDs := map[SubscriptionTier]string{
		TierStarter:    "price_starter_monthly",    // Replace with actual Stripe price IDs
		TierPro:        "price_pro_monthly",        // Replace with actual Stripe price IDs
		TierEnterprise: "price_enterprise_monthly", // Replace with actual Stripe price IDs
	}
	return priceIDs[tier]
}

// HandleWebhook processes Stripe webhook events
func (s *BillingService) HandleWebhook(event stripe.Event) error {
	switch event.Type {
	case "checkout.session.completed":
		return s.handleCheckoutCompleted(event)
	case "customer.subscription.updated":
		return s.handleSubscriptionUpdated(event)
	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(event)
	case "invoice.payment_failed":
		return s.handlePaymentFailed(event)
	default:
		// Unhandled event type
		return nil
	}
}

// handleCheckoutCompleted handles successful checkout
func (s *BillingService) handleCheckoutCompleted(event stripe.Event) error {
	var checkoutSession stripe.CheckoutSession
	if err := stripe.ParseEvent(event, &checkoutSession); err != nil {
		return err
	}

	// Find pending subscription by session ID
	var sub models.UserSubscription
	if err := s.db.Where("stripe_session_id = ?", checkoutSession.ID).First(&sub).Error; err != nil {
		return err
	}

	// Update subscription with Stripe data
	sub.StripeSubscriptionID = checkoutSession.Subscription.ID
	sub.StripeCustomerID = checkoutSession.Customer.ID
	sub.Status = "active"
	sub.CurrentPeriodStart = time.Unix(checkoutSession.Subscription.CurrentPeriodStart, 0)
	sub.CurrentPeriodEnd = time.Unix(checkoutSession.Subscription.CurrentPeriodEnd, 0)

	return s.db.Save(&sub).Error
}

// handleSubscriptionUpdated handles subscription updates
func (s *BillingService) handleSubscriptionUpdated(event stripe.Event) error {
	var stripeSub stripe.Subscription
	if err := stripe.ParseEvent(event, &stripeSub); err != nil {
		return err
	}

	// Find subscription by Stripe ID
	var sub models.UserSubscription
	if err := s.db.Where("stripe_subscription_id = ?", stripeSub.ID).First(&sub).Error; err != nil {
		return err
	}

	// Update subscription status
	sub.Status = string(stripeSub.Status)
	sub.CurrentPeriodStart = time.Unix(stripeSub.CurrentPeriodStart, 0)
	sub.CurrentPeriodEnd = time.Unix(stripeSub.CurrentPeriodEnd, 0)

	// Update tier if changed
	if tier, ok := stripeSub.Metadata["tier"]; ok {
		sub.Tier = tier
		limits := GetTierLimits(SubscriptionTier(tier))
		sub.AgentsLimit = limits.AgentsLimit
		sub.OrganizationsLimit = limits.OrgsLimit
	}

	return s.db.Save(&sub).Error
}

// handleSubscriptionDeleted handles subscription cancellation
func (s *BillingService) handleSubscriptionDeleted(event stripe.Event) error {
	var stripeSub stripe.Subscription
	if err := stripe.ParseEvent(event, &stripeSub); err != nil {
		return err
	}

	// Find subscription by Stripe ID
	var sub models.UserSubscription
	if err := s.db.Where("stripe_subscription_id = ?", stripeSub.ID).First(&sub).Error; err != nil {
		return err
	}

	// Downgrade to free tier
	sub.Status = "canceled"
	sub.Tier = string(TierFree)
	sub.AgentsLimit = 3
	sub.OrganizationsLimit = 1

	return s.db.Save(&sub).Error
}

// handlePaymentFailed handles failed payments
func (s *BillingService) handlePaymentFailed(event stripe.Event) error {
	var invoice stripe.Invoice
	if err := stripe.ParseEvent(event, &invoice); err != nil {
		return err
	}

	// Find subscription by Stripe ID
	var sub models.UserSubscription
	if err := s.db.Where("stripe_subscription_id = ?", invoice.Subscription.ID).First(&sub).Error; err != nil {
		return err
	}

	// Mark subscription as past_due
	sub.Status = "past_due"

	return s.db.Save(&sub).Error
}

// GetUserSubscription retrieves a user's subscription
func (s *BillingService) GetUserSubscription(userID string) (*models.UserSubscription, error) {
	var sub models.UserSubscription
	err := s.db.Where("user_id = ? AND status IN ?", userID, []string{"active", "past_due"}).First(&sub).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return free tier if no subscription
			return &models.UserSubscription{
				UserID:              userID,
				Tier:                string(TierFree),
				Status:              "active",
				AgentsLimit:         3,
				OrganizationsLimit:  1,
			}, nil
		}
		return nil, err
	}

	return &sub, nil
}

// CheckAgentLimit checks if user can create more agents
func (s *BillingService) CheckAgentLimit(userID string) error {
	sub, err := s.GetUserSubscription(userID)
	if err != nil {
		return err
	}

	if sub.Status != "active" {
		return errors.New("subscription is not active")
	}

	// Count current agents
	var agentCount int64
	if err := s.db.Model(&models.Agent{}).Where("user_id = ?", userID).Count(&agentCount).Error; err != nil {
		return err
	}

	// Check limit (unlimited = -1)
	if sub.AgentsLimit > 0 && int(agentCount) >= sub.AgentsLimit {
		return errors.New("agent limit reached for current subscription tier. Please upgrade to create more agents")
	}

	return nil
}

// CheckOrganizationLimit checks if user can create more organizations
func (s *BillingService) CheckOrganizationLimit(userID string) error {
	sub, err := s.GetUserSubscription(userID)
	if err != nil {
		return err
	}

	if sub.Status != "active" {
		return errors.New("subscription is not active")
	}

	// Count current organizations
	var orgCount int64
	if err := s.db.Model(&models.Organisation{}).Where("owner_id = ?", userID).Count(&orgCount).Error; err != nil {
		return err
	}

	// Check limit (unlimited = -1)
	if sub.OrganizationsLimit > 0 && int(orgCount) >= sub.OrganizationsLimit {
		return errors.New("organization limit reached for current subscription tier. Please upgrade to create more organizations")
	}

	return nil
}

// CancelSubscription cancels a user's subscription
func (s *BillingService) CancelSubscription(userID string) error {
	sub, err := s.GetUserSubscription(userID)
	if err != nil {
		return err
	}

	// Cancel in Stripe
	if sub.StripeSubscriptionID != "" {
		_, err = subscription.Cancel(sub.StripeSubscriptionID, nil)
		if err != nil {
			return err
		}
	}

	// Update local subscription
	sub.Status = "canceled"
	sub.Tier = string(TierFree)
	sub.AgentsLimit = 3
	sub.OrganizationsLimit = 1

	return s.db.Save(&sub).Error
}

// UpdateSubscription updates a subscription tier
func (s *BillingService) UpdateSubscription(userID string, newTier SubscriptionTier) error {
	sub, err := s.GetUserSubscription(userID)
	if err != nil {
		return err
	}

	// For free tier, just update local
	if newTier == TierFree {
		sub.Tier = string(TierFree)
		sub.AgentsLimit = 3
		sub.OrganizationsLimit = 1
		return s.db.Save(&sub).Error
	}

	// For paid tiers, create new checkout session
	checkoutSession, err := s.createCheckoutSession(sub.UserID, newTier)
	if err != nil {
		return err
	}

	// Update subscription
	sub.StripeSessionID = checkoutSession.ID
	sub.Status = "pending"

	return s.db.Save(&sub).Error
}
