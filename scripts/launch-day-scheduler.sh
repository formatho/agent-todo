#!/bin/bash

# Product Hunt Launch Day Social Media Scheduler
# Automated posting for Agent-Todo Product Hunt launch

LAUNCH_DATE="2026-03-30"
PLATFORMS=("twitter" "linkedin" "devto")
TIMEZONE="Asia/Calcutta"

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Agent-Todo Product Hunt Launch Day Scheduler${NC}"
echo "Launch Date: $LAUNCH_DATE"
echo "Timezone: $TIMEZONE"
echo "====================================="

# Function to check if launch day has arrived
check_launch_day() {
    current_date=$(date +%Y-%m-%d)
    if [[ "$current_date" == "$LAUNCH_DATE" ]]; then
        return 0
    else
        echo -e "${YELLOW}⏰ Launch day is not yet here. Current date: $current_date${NC}"
        return 1
    fi
}

# Function to get platform posting times
get_posting_schedule() {
    case $1 in
        "twitter")
            echo "9:00,10:00,11:00,12:00,13:00,14:00,15:00,16:00,17:00,18:00"
            ;;
        "linkedin")
            echo "9:00,11:00,13:00,15:00,17:00"
            ;;
        "devto")
            echo "9:00"
            ;;
        "reddit")
            echo "10:00"
            ;;
        "facebook")
            echo "12:00"
            ;;
        *)
            echo ""
            ;;
    esac
}

# Function to get platform post count
get_post_count() {
    case $1 in
        "twitter")
            echo "10"
            ;;
        "linkedin")
            echo "5"
            ;;
        "devto")
            echo "1"
            ;;
        "reddit")
            echo "1"
            ;;
        "facebook")
            echo "1"
            ;;
        *)
            echo "0"
            ;;
    esac
}

# Function to show posting schedule
show_schedule() {
    echo -e "\n${GREEN}📅 Launch Day Posting Schedule${NC}"
    echo "====================================="
    
    for platform in "${PLATFORMS[@]}"; do
        times=$(get_posting_schedule "$platform")
        count=$(get_post_count "$platform")
        echo -e "${BLUE}$platform:${NC} $count posts at $times"
    done
    
    echo -e "\n${YELLOW}Additional Platforms:${NC}"
    echo "Reddit: 1 post at 10:00"
    echo "Facebook: 1 post at 12:00"
}

# Function to validate social media content
validate_content() {
    echo -e "\n${GREEN}📝 Validating Social Media Content${NC}"
    echo "====================================="
    
    content_file="/Users/studio/sandbox/formatho/agent-todo/PRODUCT_HUNT_SOCIAL_MEDIA.md"
    
    if [[ ! -f "$content_file" ]]; then
        echo -e "${RED}❌ Content file not found: $content_file${NC}"
        return 1
    fi
    
    # Check word count
    word_count=$(wc -w < "$content_file")
    echo -e "✅ Content file found: $word_count words"
    
    # Check for required sections
    if grep -q "Twitter/X Content" "$content_file"; then
        echo -e "✅ Twitter content included"
    else
        echo -e "${RED}❌ Twitter content missing${NC}"
    fi
    
    if grep -q "LinkedIn Content" "$content_file"; then
        echo -e "✅ LinkedIn content included"
    else
        echo -e "${RED}❌ LinkedIn content missing${NC}"
    fi
    
    if grep -q "Posting Schedule" "$content_file"; then
        echo -e "✅ Posting schedule included"
    else
        echo -e "${RED}❌ Posting schedule missing${NC}"
    fi
    
    echo -e "${GREEN}✅ Content validation complete${NC}"
}

# Function to check launch readiness
check_readiness() {
    echo -e "\n${GREEN}🚀 Launch Readiness Checklist${NC}"
    echo "====================================="
    
    # Check if website is accessible
    if curl -s -o /dev/null -w "%{http_code}" https://todo.formatho.com | grep -q "200\|302"; then
        echo -e "✅ Website is accessible"
    else
        echo -e "${RED}❌ Website is not accessible${NC}"
    fi
    
    # Check if Product Hunt listing exists
    if curl -s -o /dev/null -w "%{http_code}" https://producthunt.com/products/agent-todo | grep -q "200"; then
        echo -e "✅ Product Hunt listing exists"
    else
        echo -e "${YELLOW}⚠️  Product Hunt listing may not be live yet${NC}"
    fi
    
    # Check monitoring system
    if [[ -f "/Users/studio/sandbox/formatho/agent-todo/backend/bin/agent-todo" ]]; then
        echo -e "✅ Backend monitoring system compiled"
    else
        echo -e "${YELLOW}⚠️  Backend system may need compilation${NC}"
    fi
    
    # Check social media content
    if [[ -f "/Users/studio/sandbox/formatho/agent-todo/PRODUCT_HUNT_SOCIAL_MEDIA.md" ]]; then
        echo -e "✅ Social media content prepared"
    else
        echo -e "${RED}❌ Social media content not found${NC}"
    fi
    
    echo -e "\n${GREEN}✅ Launch readiness check complete${NC}"
}

# Function to show countdown
show_countdown() {
    if check_launch_day; then
        echo -e "\n${GREEN}🎉 LAUNCH DAY IS HERE! 🎉${NC}"
        echo "====================================="
        echo "Time to activate the social media machine! 🚀"
    else
        # Calculate days until launch
        launch_date=$(date -d "$LAUNCH_DATE" +%s)
        current_date=$(date +%s)
        days_left=$(( (launch_date - current_date) / 86400 ))
        
        echo -e "\n${YELLOW}⏰ Countdown to Launch Day${NC}"
        echo "====================================="
        echo "Days remaining: $days_left"
        
        if [[ $days_left -le 7 ]]; then
            echo -e "${YELLOW}Final preparations phase - review all assets!${NC}"
        elif [[ $days_left -le 14 ]]; then
            echo -e "${BLUE}Content creation phase - prepare social posts${NC}"
        else
            echo -e "${GREEN}Early planning phase - build strategy${NC}"
        fi
    fi
}

# Function to generate posting reminders
generate_reminders() {
    echo -e "\n${GREEN}📱 Posting Reminders${NC}"
    echo "====================================="
    
    if check_launch_day; then
        echo "TODAY - March 30, 2026:"
        echo "• 9:00 AM: Launch announcement on all platforms"
        echo "• 10:00 AM: Feature spotlight (Twitter)"
        echo "• 11:00 AM: Problem/Solution (Twitter)"
        echo "• 11:00 AM: Industry insight (LinkedIn)"
        echo "• 12:00 PM: Use case (Twitter)"
        echo "• 12:00 PM: Facebook post"
        echo "• 1:00 PM: Behind the scenes (Twitter)"
        echo "• 1:00 PM: Case study preview (LinkedIn)"
        echo "• 2:00 PM: Interactive poll (Twitter)"
        echo "• 3:00 PM: Thought leadership (LinkedIn)"
        echo "• 4:00 PM: Call to action (Twitter)"
        echo "• 5:00 PM: Evening wrap-up (Twitter)"
        echo "• 5:00 PM: Future vision (LinkedIn)"
        
        echo -e "\n${YELLOW}⚠️  Don't forget to:${NC}"
        echo "• Respond to comments within 1 hour"
        echo "• Retweet and share user content"
        echo "• Follow back Product Hunt voters"
        echo "• Share live updates for milestones"
        echo "• Thank users for their support"
    else
        echo "Pre-launch reminders:"
        echo "• Review all social media content"
        echo "• Test website accessibility"
        echo "• Verify Product Hunt listing"
        echo "• Prepare monitoring dashboards"
        echo "• Plan team response strategy"
    fi
}

# Function to show analytics tracking
show_analytics() {
    echo -e "\n${GREEN}📊 Analytics Tracking${NC}"
    echo "====================================="
    
    echo "Track these metrics on launch day:"
    echo ""
    echo "Twitter/X:"
    echo "• Impressions target: 500+"
    echo "• Engagement target: 50+"
    echo "• Follower growth: 20+"
    echo ""
    echo "LinkedIn:"
    echo "• Impressions target: 200+"
    echo "• Engagement target: 25+"
    echo "• Profile visits: 50+"
    echo ""
    echo "Product Hunt:"
    echo "• Votes target: 100+"
    echo "• Comments target: 50+"
    echo "• Upvotes target: 200+"
    echo ""
    echo "Website:"
    echo "• Traffic target: 500+ unique visitors"
    echo "• Sign-ups target: 50+"
    echo "• Bounce rate target: < 50%"
}

# Main menu
show_menu() {
    echo -e "\n${BLUE}🎯 What would you like to do?${NC}"
    echo "====================================="
    echo "1. Show launch day posting schedule"
    echo "2. Validate social media content"
    echo "3. Check launch readiness"
    echo "4. Show countdown to launch"
    echo "5. Generate posting reminders"
    echo "6. Show analytics targets"
    echo "7. Exit"
    
    read -p "Enter your choice (1-7): " choice
    
    case $choice in
        1)
            show_schedule
            ;;
        2)
            validate_content
            ;;
        3)
            check_readiness
            ;;
        4)
            show_countdown
            ;;
        5)
            generate_reminders
            ;;
        6)
            show_analytics
            ;;
        7)
            echo -e "${GREEN}👋 Good luck with the launch! 🚀${NC}"
            exit 0
            ;;
        *)
            echo -e "${RED}❌ Invalid choice. Please try again.${NC}"
            ;;
    esac
}

# Main execution
while true; do
    show_menu
    echo -e "\nPress Enter to continue..."
    read
done