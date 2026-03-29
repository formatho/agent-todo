#!/bin/bash

# Reddit API Setup Script for Agent-Todo
# This script helps set up Reddit API access for content distribution

echo "🚀 Setting up Reddit API for Agent-Todo Content Distribution"
echo "=========================================================="

# Check if Reddit credentials are already configured
if [ -f .env ]; then
    echo "✅ .env file exists"
    if grep -q "REDDIT_CLIENT_ID" .env; then
        echo "✅ Reddit credentials already configured"
        echo "   REDDIT_CLIENT_ID: $(grep REDDIT_CLIENT_ID .env | cut -d'=' -f2)"
    else
        echo "❌ Reddit credentials not found in .env file"
    fi
else
    echo "❌ .env file not found, copying from .env.example"
    cp .env.example .env
fi

echo ""
echo "📋 Reddit API Setup Instructions:"
echo "==============================="
echo ""
echo "1. Create a Reddit App:"
echo "   - Go to: https://www.reddit.com/prefs/apps"
echo "   - Click 'are you a developer? create an app'"
echo "   - Choose 'script' (for personal use)"
echo "   - Name: 'Agent-Todo Content Distribution'"
echo "   - About URL: https://todo.formatho.com"
echo "   - Redirect URI: http://localhost:8080 (can be anything for script)"
echo ""
echo "2. Get your credentials:"
echo "   - Copy the 'personal use script' Client ID"
echo "   - Copy the Client Secret"
echo "   - You'll need your Reddit username and password"
echo ""
echo "3. Add to your .env file:"
echo "   REDDIT_CLIENT_ID=your-client-id-here"
echo "   REDDIT_CLIENT_SECRET=your-client-secret-here"
echo "   REDDIT_USERNAME=your-reddit-username"
echo "   REDDIT_PASSWORD=your-reddit-password"
echo "   REDDIT_USER_AGENT=Agent-Todo/1.0"
echo ""
echo "4. Test Reddit posting:"
echo "   cd scripts"
echo "   chmod +x reddit-poster"
echo "   ./reddit-poster ../blog/building-persistent-memory-for-ai-agents.md"
echo ""
echo "⚠️  Important Notes:"
echo "   - Reddit requires account age and karma for posting"
echo "   - New accounts may be restricted from posting"
echo "   - Use the 'flair_text' parameter for better visibility"
echo "   - Test with a simple post first before automated distribution"
echo ""
echo "🔧 Alternative: Manual Posting"
echo "==============================="
echo "If automated posting fails, you can manually post to:"
echo "   r/artificial - AI tools and discussions"
echo "   r/MachineLearning - ML and AI content"
echo "   r/devops - DevOps and automation"
echo "   r/programming - Programming tutorials"
echo ""
echo "Include these tags: #AI #TaskManagement #AgentTodo #AIOps"
echo ""
echo "🎯 Next Steps:"
echo "============="
echo "1. Set up Reddit API credentials"
echo "2. Test with a sample blog post"
echo "3. Schedule automated posting times"
echo "4. Monitor and adjust based on engagement"
echo ""