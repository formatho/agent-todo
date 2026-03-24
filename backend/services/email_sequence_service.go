package services

import (
	"fmt"
	"time"

	"github.com/formatho/agent-todo/db"
	"github.com/formatho/agent-todo/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// EmailSequenceService handles email sequence management
type EmailSequenceService struct {
	db           *gorm.DB
	emailService *EmailService
}

// NewEmailSequenceService creates a new email sequence service
func NewEmailSequenceService(db *gorm.DB) *EmailSequenceService {
	return &EmailSequenceService{
		db:           db,
		emailService: NewEmailService(db),
	}
}

// SeedTrialConversionSequence creates the default trial-to-paid conversion sequence
func (s *EmailSequenceService) SeedTrialConversionSequence() error {
	// Check if sequence already exists
	var existing models.EmailSequence
	if err := s.db.Where("trigger = ?", "trial_started").First(&existing).Error; err == nil {
		return nil // Already exists
	}

	// Create email templates
	templates := []models.EmailTemplate{
		{
			Name:         "Trial Welcome - Day 1",
			Subject:      "Welcome to Agent Todo! Your 14-Day Trial Has Started 🚀",
			TemplateType: models.EmailTemplateTypeWelcome,
			BodyHTML: `<h1>Welcome to Agent Todo!</h1>
<p>Hi {{.UserName}},</p>
<p>Congratulations on starting your 14-day free trial of Agent Todo! You now have access to powerful AI-powered task management.</p>

<h2>Here's what you can do:</h2>
<ul>
  <li>✅ Create and assign tasks to AI agents</li>
  <li>✅ Organize projects with your team</li>
  <li>✅ Track progress in real-time</li>
  <li>✅ Integrate with your existing tools</li>
</ul>

<h2>Quick Start Guide:</h2>
<ol>
  <li>Create your first project</li>
  <li>Add your AI agents</li>
  <li>Create tasks and watch them get done!</li>
</ol>

<p><a href="{{.DashboardURL}}" style="background: #4F46E5; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px;">Go to Dashboard</a></p>

<p>If you have any questions, just reply to this email.</p>

<p>Cheers,<br>The Formatho Team</p>`,
			BodyText: `Welcome to Agent Todo!

Hi {{.UserName}},

Congratulations on starting your 14-day free trial of Agent Todo! You now have access to powerful AI-powered task management.

Here's what you can do:
- Create and assign tasks to AI agents
- Organize projects with your team
- Track progress in real-time
- Integrate with your existing tools

Quick Start Guide:
1. Create your first project
2. Add your AI agents
3. Create tasks and watch them get done!

Go to Dashboard: {{.DashboardURL}}

If you have any questions, just reply to this email.

Cheers,
The Formatho Team`,
			Variables: `["UserName", "DashboardURL"]`,
			IsActive:  true,
		},
		{
			Name:         "Value Tips - Day 3",
			Subject:      "3 Pro Tips to Supercharge Your AI Agents ⚡",
			TemplateType: models.EmailTemplateTypeValueTips,
			BodyHTML: `<h1>Unlock More Value from Agent Todo</h1>
<p>Hi {{.UserName}},</p>
<p>You've been using Agent Todo for a few days now. Here are 3 pro tips to help you get even more done:</p>

<h2>💡 Tip 1: Use Project Context</h2>
<p>Add LLM context to your projects so agents understand your coding standards, deployment URLs, and documentation links. This helps them work more autonomously.</p>

<h2>💡 Tip 2: Assign Tasks by Priority</h2>
<p>Use priority levels (low, medium, high, critical) to help agents focus on what matters most. Critical tasks always get picked up first.</p>

<h2>💡 Tip 3: Add Comments for Context</h2>
<p>When you need to guide an agent, add comments to tasks. Agents can read comments and adjust their approach accordingly.</p>

<p><a href="{{.DashboardURL}}" style="background: #4F46E5; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px;">Try These Tips Now</a></p>

<p>Your trial has {{.DaysRemaining}} days left. Make the most of it!</p>

<p>Cheers,<br>The Formatho Team</p>`,
			BodyText: `Unlock More Value from Agent Todo

Hi {{.UserName}},

You've been using Agent Todo for a few days now. Here are 3 pro tips to help you get even more done:

💡 Tip 1: Use Project Context
Add LLM context to your projects so agents understand your coding standards, deployment URLs, and documentation links. This helps them work more autonomously.

💡 Tip 2: Assign Tasks by Priority
Use priority levels (low, medium, high, critical) to help agents focus on what matters most. Critical tasks always get picked up first.

💡 Tip 3: Add Comments for Context
When you need to guide an agent, add comments to tasks. Agents can read comments and adjust their approach accordingly.

Try These Tips Now: {{.DashboardURL}}

Your trial has {{.DaysRemaining}} days left. Make the most of it!

Cheers,
The Formatho Team`,
			Variables: `["UserName", "DashboardURL", "DaysRemaining"]`,
			IsActive:  true,
		},
		{
			Name:         "Case Study - Day 7",
			Subject:      "How One Team Saved 20 Hours/Week with Agent Todo 📊",
			TemplateType: models.EmailTemplateTypeCaseStudy,
			BodyHTML: `<h1>Real Results from Real Teams</h1>
<p>Hi {{.UserName}},</p>
<p>Halfway through your trial! Let us share how others are using Agent Todo:</p>

<h2>Case Study: DevOps Team at StreamFlow</h2>
<blockquote style="border-left: 4px solid #4F46E5; padding-left: 16px; margin: 16px 0;">
<p>"Before Agent Todo, our 5-person DevOps team spent 20+ hours/week on routine maintenance tasks—updating dependencies, rotating credentials, monitoring alerts. Now our AI agents handle 80% of these tasks autonomously. We ship faster and our engineers focus on architecture instead of toil."</p>
<footer><strong>— Sarah Chen, VP of Engineering</strong></footer>
</blockquote>

<h2>What StreamFlow Automated:</h2>
<ul>
  <li>✅ Dependency updates with auto-PR creation</li>
  <li>✅ SSL certificate renewal monitoring</li>
  <li>✅ Log rotation and cleanup</li>
  <li>✅ On-call alert triage and escalation</li>
</ul>

<h2>Your Results Could Be Next</h2>
<p>Think about repetitive tasks in your workflow. What could you automate?</p>

<p><a href="{{.DashboardURL}}" style="background: #4F46E5; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px;">Create Your First Automation</a></p>

<p>{{.DaysRemaining}} days left in your trial.</p>

<p>Cheers,<br>The Formatho Team</p>`,
			BodyText: `Real Results from Real Teams

Hi {{.UserName}},

Halfway through your trial! Let us share how others are using Agent Todo:

Case Study: DevOps Team at StreamFlow
"Before Agent Todo, our 5-person DevOps team spent 20+ hours/week on routine maintenance tasks—updating dependencies, rotating credentials, monitoring alerts. Now our AI agents handle 80% of these tasks autonomously. We ship faster and our engineers focus on architecture instead of toil."
— Sarah Chen, VP of Engineering

What StreamFlow Automated:
- Dependency updates with auto-PR creation
- SSL certificate renewal monitoring
- Log rotation and cleanup
- On-call alert triage and escalation

Your Results Could Be Next
Think about repetitive tasks in your workflow. What could you automate?

Create Your First Automation: {{.DashboardURL}}

{{.DaysRemaining}} days left in your trial.

Cheers,
The Formatho Team`,
			Variables: `["UserName", "DashboardURL", "DaysRemaining"]`,
			IsActive:  true,
		},
		{
			Name:         "Limited Offer - Day 12",
			Subject:      "🎁 Special Offer: 20% Off Your First Year",
			TemplateType: models.EmailTemplateTypeLimitedOffer,
			BodyHTML: `<h1>An Exclusive Offer for You</h1>
<p>Hi {{.UserName}},</p>
<p>Your trial ends in just <strong>2 days</strong>. We've loved having you, and we want to make your decision easy.</p>

<div style="background: linear-gradient(135deg, #4F46E5 0%, #7C3AED 100%); color: white; padding: 24px; border-radius: 12px; margin: 24px 0;">
  <h2 style="margin-top: 0; color: white;">🎁 20% Off Your First Year</h2>
  <p style="font-size: 18px;">Use code <strong style="font-size: 24px;">{{.DiscountCode}}</strong> at checkout</p>
  <p style="margin-bottom: 0;">Valid until {{.OfferExpires}}</p>
</div>

<h2>What You'll Keep:</h2>
<ul>
  <li>✅ Unlimited AI agents</li>
  <li>✅ All your tasks and projects</li>
  <li>✅ Priority support</li>
  <li>✅ Advanced analytics</li>
</ul>

<p>This offer is only available for the next 48 hours. Don't miss out!</p>

<p><a href="{{.PricingURL}}" style="background: #4F46E5; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px;">Claim Your Discount</a></p>

<p>Cheers,<br>The Formatho Team</p>`,
			BodyText: `An Exclusive Offer for You

Hi {{.UserName}},

Your trial ends in just 2 days. We've loved having you, and we want to make your decision easy.

🎁 20% OFF YOUR FIRST YEAR
Use code {{.DiscountCode}} at checkout
Valid until {{.OfferExpires}}

What You'll Keep:
- Unlimited AI agents
- All your tasks and projects
- Priority support
- Advanced analytics

This offer is only available for the next 48 hours. Don't miss out!

Claim Your Discount: {{.PricingURL}}

Cheers,
The Formatho Team`,
			Variables: `["UserName", "PricingURL", "DiscountCode", "OfferExpires"]`,
			IsActive:  true,
		},
		{
			Name:         "Upgrade Reminder - Day 14",
			Subject:      "Your Trial Ends Today ⏰ Don't Lose Access",
			TemplateType: models.EmailTemplateTypeUpgradeReminder,
			BodyHTML: `<h1>Final Reminder: Your Trial Ends Today</h1>
<p>Hi {{.UserName}},</p>
<p>This is it—your 14-day trial of Agent Todo ends <strong>today</strong>.</p>

<h2>What Happens Tomorrow:</h2>
<ul>
  <li>❌ Your AI agents will be paused</li>
  <li>❌ Tasks will become read-only</li>
  <li>❌ You'll lose access to premium features</li>
</ul>

<h2>Keep Everything Running:</h2>
<p>Upgrade now and your agents will keep working without interruption.</p>

<div style="background: #FEF3C7; padding: 16px; border-radius: 8px; margin: 16px 0;">
  <p style="margin: 0;">💡 <strong>Remember:</strong> Use code <strong>{{.DiscountCode}}</strong> for 20% off your first year!</p>
</div>

<p><a href="{{.PricingURL}}" style="background: #4F46E5; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; font-weight: bold;">Upgrade Now</a></p>

<p>Questions? Just reply to this email—we're here to help.</p>

<p>Thanks for trying Agent Todo!<br>The Formatho Team</p>`,
			BodyText: `Final Reminder: Your Trial Ends Today

Hi {{.UserName}},

This is it—your 14-day trial of Agent Todo ends today.

What Happens Tomorrow:
- Your AI agents will be paused
- Tasks will become read-only
- You'll lose access to premium features

Keep Everything Running:
Upgrade now and your agents will keep working without interruption.

💡 Remember: Use code {{.DiscountCode}} for 20% off your first year!

Upgrade Now: {{.PricingURL}}

Questions? Just reply to this email—we're here to help.

Thanks for trying Agent Todo!
The Formatho Team`,
			Variables: `["UserName", "PricingURL", "DiscountCode"]`,
			IsActive:  true,
		},
	}

	// Create templates
	for i := range templates {
		if err := s.db.Create(&templates[i]).Error; err != nil {
			return fmt.Errorf("failed to create template %s: %w", templates[i].Name, err)
		}
	}

	// Create sequence
	sequence := &models.EmailSequence{
		Name:        "Trial to Paid Conversion",
		Description: "5-email nurture sequence to convert trial users to paid subscribers",
		Trigger:     "trial_started",
		Status:      models.EmailSequenceStatusActive,
	}

	if err := s.db.Create(sequence).Error; err != nil {
		return fmt.Errorf("failed to create sequence: %w", err)
	}

	// Create sequence steps
	steps := []models.EmailSequenceStep{
		{SequenceID: sequence.ID, TemplateID: templates[0].ID, DelayDays: 1, Order: 1},   // Day 1: Welcome
		{SequenceID: sequence.ID, TemplateID: templates[1].ID, DelayDays: 3, Order: 2},   // Day 3: Value Tips
		{SequenceID: sequence.ID, TemplateID: templates[2].ID, DelayDays: 7, Order: 3},   // Day 7: Case Study
		{SequenceID: sequence.ID, TemplateID: templates[3].ID, DelayDays: 12, Order: 4},  // Day 12: Limited Offer
		{SequenceID: sequence.ID, TemplateID: templates[4].ID, DelayDays: 14, Order: 5},  // Day 14: Final Reminder
	}

	for i := range steps {
		if err := s.db.Create(&steps[i]).Error; err != nil {
			return fmt.Errorf("failed to create sequence step: %w", err)
		}
	}

	return nil
}

// TriggerSequence starts an email sequence for a user
func (s *EmailSequenceService) TriggerSequence(trigger string, userID, organisationID uuid.UUID) error {
	// Find the sequence
	var sequence models.EmailSequence
	if err := s.db.Preload("Steps.Template").Where("trigger = ? AND status = ?", trigger, models.EmailSequenceStatusActive).First(&sequence).Error; err != nil {
		return fmt.Errorf("sequence not found: %w", err)
	}

	// Get user
	var user models.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Queue each step
	for _, step := range sequence.Steps {
		vars := s.getTemplateVariables(&user, &organisationID, step.Template.TemplateType, step.DelayDays)

		bodyHTML, err := s.emailService.RenderTemplate(step.Template.BodyHTML, vars)
		if err != nil {
			return fmt.Errorf("failed to render HTML template: %w", err)
		}

		bodyText, err := s.emailService.RenderTemplate(step.Template.BodyText, vars)
		if err != nil {
			return fmt.Errorf("failed to render text template: %w", err)
		}

		subject, err := s.emailService.RenderTemplate(step.Template.Subject, vars)
		if err != nil {
			return fmt.Errorf("failed to render subject: %w", err)
		}

		sequenceIDStr := sequence.ID.String()
		stepIDStr := step.ID.String()

		_, err = s.emailService.QueueSequenceEmail(
			userID.String(),
			organisationID.String(),
			&sequenceIDStr,
			&stepIDStr,
			user.Email,
			subject,
			bodyHTML,
			bodyText,
			step.DelayDays,
		)
		if err != nil {
			return fmt.Errorf("failed to queue email: %w", err)
		}
	}

	return nil
}

// getTemplateVariables returns template variables based on email type
func (s *EmailSequenceService) getTemplateVariables(user *models.User, organisationID *uuid.UUID, templateType models.EmailTemplateType, delayDays int) map[string]interface{} {
	baseURL := "https://todo.formatho.com"
	vars := map[string]interface{}{
		"UserName":      user.Email, // Default to email, can be improved with name field
		"DashboardURL":  fmt.Sprintf("%s/dashboard", baseURL),
		"PricingURL":    fmt.Sprintf("%s/pricing", baseURL),
		"DiscountCode":  "TRIAL20",
		"DaysRemaining": 14 - delayDays,
		"OfferExpires":  time.Now().AddDate(0, 0, 7).Format("January 2, 2006"),
	}

	// Calculate days remaining more accurately based on trial end date
	var subscription models.Subscription
	if err := s.db.Where("organisation_id = ?", organisationID).First(&subscription).Error; err == nil {
		if subscription.TrialEndsAt != nil {
			daysRemaining := int(time.Until(*subscription.TrialEndsAt).Hours() / 24)
			if daysRemaining > 0 {
				vars["DaysRemaining"] = daysRemaining
			}
		}
	}

	return vars
}

// ProcessEmailQueue processes pending emails (called by cron)
func (s *EmailSequenceService) ProcessEmailQueue(batchSize int) (int, int, error) {
	return s.emailService.ProcessQueue(batchSize)
}

// EnsureSequenceTables creates the necessary database tables
func (s *EmailSequenceService) EnsureSequenceTables() error {
	return s.db.AutoMigrate(
		&models.Subscription{},
		&models.EmailTemplate{},
		&models.EmailSequence{},
		&models.EmailSequenceStep{},
		&models.EmailQueue{},
		&models.EmailLog{},
	)
}

// StartTrial starts a trial for an organisation
func (s *EmailSequenceService) StartTrial(organisationID uuid.UUID) error {
	now := time.Now()
	trialEnd := now.AddDate(0, 0, 14)

	subscription := &models.Subscription{
		OrganisationID: organisationID,
		Plan:           models.SubscriptionPlanTrial,
		Status:         models.SubscriptionStatusTrialing,
		TrialStartsAt:  &now,
		TrialEndsAt:    &trialEnd,
	}

	if err := s.db.Create(subscription).Error; err != nil {
		return err
	}

	// Get organisation owner
	var org models.Organisation
	if err := s.db.Preload("CreatedBy").First(&org, "id = ?", organisationID).Error; err != nil {
		return err
	}

	// Trigger the email sequence
	return s.TriggerSequence("trial_started", org.CreatedByUserID, organisationID)
}

// GetSubscription returns subscription for an organisation
func (s *EmailSequenceService) GetSubscription(organisationID uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription
	if err := s.db.Where("organisation_id = ?", organisationID).First(&subscription).Error; err != nil {
		return nil, err
	}
	return &subscription, nil
}

// CancelTrial cancels an ongoing trial
func (s *EmailSequenceService) CancelTrial(organisationID uuid.UUID) error {
	now := time.Now()
	return s.db.Model(&models.Subscription{}).
		Where("organisation_id = ? AND status = ?", organisationID, models.SubscriptionStatusTrialing).
		Updates(map[string]interface{}{
			"status":     models.SubscriptionStatusCanceled,
			"canceled_at": now,
		}).Error
}
