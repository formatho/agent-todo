# OpenClaw Skills Review

**Date:** March 8, 2026  
**Reviewer:** Agent-Todo  
**Purpose:** Comprehensive review of all available OpenClaw skills, their functionality, capabilities, and integration points.

---

## Overview

OpenClaw skills are modular capabilities that extend agent functionality. Each skill provides specialized instructions and tools for specific tasks. This document catalogs all available skills and their potential use cases.

---

## Available Skills

### 1. **21st-dev** - React Website Builder
**Description:** Build premium React websites from scratch with award-winning design patterns.

**Capabilities:**
- Generates landing pages, dashboards, auth flows, and marketing sites
- 61 component blueprints with Framer Motion
- Glassmorphism and micro-interactions
- Modern UI/UX patterns

**Use Cases:**
- Rapid prototyping of web applications
- Creating marketing landing pages
- Building admin dashboards
- Developing authentication flows

**Integration Points:** Frontend development workflows

---

### 2. **AI Trend Hunter Pro** - Automated Content Agent
**Description:** High-performance automation for viral social media content generation.

**Capabilities:**
- Scrapes global trends
- Generates posts for X (Twitter), Xiaohongshu, and LinkedIn
- Based on tech influencer methodology

**Use Cases:**
- Social media automation
- Content creation for multiple platforms
- Trend monitoring and content generation

**Integration Points:** Social media APIs, content scheduling platforms

---

### 3. **coding-agent** - Development Task Delegation
**Description:** Delegate coding tasks to specialized AI agents (Codex, Claude Code, Pi).

**Capabilities:**
- Background process execution
- PR reviews
- Large codebase refactoring
- Iterative coding with file exploration

**Use Cases:**
- Complex feature development
- Code review automation
- Refactoring legacy code
- Multi-file development tasks

**Integration Points:** GitHub workflows, CI/CD pipelines

**Limitations:** Not for simple one-liner fixes, code reading, or ~/clawd workspace

---

### 4. **gemini** - Quick Q&A and Generation
**Description:** Gemini CLI for one-shot queries, summaries, and content generation.

**Capabilities:**
- Fast Q&A
- Text summarization
- Content generation

**Use Cases:**
- Quick research queries
- Document summarization
- Draft generation

**Integration Points:** Text processing workflows

---

### 5. **gh-issues** - GitHub Issues Automation
**Description:** Fetch GitHub issues, implement fixes, and open PRs automatically.

**Capabilities:**
- Issue fetching with filters (label, milestone, assignee)
- Sub-agent spawning for fixes
- PR creation and monitoring
- Review comment handling
- Scheduled monitoring with cron

**Use Cases:**
- Automated bug fixing
- Issue triage and resolution
- PR workflow automation
- Continuous integration maintenance

**Integration Points:** GitHub API, issue trackers

---

### 6. **github** - GitHub Operations
**Description:** GitHub CLI operations for issues, PRs, CI runs, and code review.

**Capabilities:**
- PR status checking
- CI monitoring
- Issue creation/commenting
- PR/issue filtering and listing
- Run log viewing

**Use Cases:**
- CI/CD monitoring
- Issue management
- Code review workflows
- Repository maintenance

**Integration Points:** GitHub API, `gh` CLI

---

### 7. **healthcheck** - Security Hardening
**Description:** Host security hardening and risk-tolerance configuration.

**Capabilities:**
- Security audits
- Firewall/SSH hardening
- Update management
- Risk posture assessment
- Exposure review
- Periodic health checks

**Use Cases:**
- Server security audits
- Compliance checking
- Automated security monitoring
- Risk assessment

**Integration Points:** System administration, OpenClaw cron scheduling

---

### 8. **himalaya** - Email Management
**Description:** CLI for email management via IMAP/SMTP.

**Capabilities:**
- List, read, write emails
- Reply and forward
- Search and organize
- Multiple account support
- MML (MIME Meta Language) composition

**Use Cases:**
- Email automation
- Inbox management
- Email-based workflows
- Communication automation

**Integration Points:** Email servers (IMAP/SMTP)

---

### 9. **imsg** - iMessage/SMS Management
**Description:** CLI for iMessage and SMS management via Messages.app.

**Capabilities:**
- List chats
- View history
- Send messages

**Use Cases:**
- Message automation on macOS
- Chat history analysis
- Communication workflows

**Integration Points:** macOS Messages.app

---

### 10. **mcporter** - MCP Server Management
**Description:** CLI to manage and call MCP servers/tools.

**Capabilities:**
- List, configure MCP servers
- Authentication management
- Direct tool calls (HTTP/stdio)
- Ad-hoc server management
- Config editing
- CLI/type generation

**Use Cases:**
- MCP server integration
- Tool orchestration
- API gateway management

**Integration Points:** MCP protocol, HTTP/stdio transports

---

### 11. **nano-pdf** - PDF Editing
**Description:** Edit PDFs with natural-language instructions.

**Capabilities:**
- Natural language PDF manipulation
- Document editing

**Use Cases:**
- Document processing
- PDF automation
- Report generation

**Integration Points:** Document workflows

---

### 12. **peekaboo** - macOS UI Automation
**Description:** Capture and automate macOS UI with Peekaboo CLI.

**Capabilities:**
- UI capture
- Automation scripting

**Use Cases:**
- UI testing
- Screenshot automation
- macOS workflow automation

**Integration Points:** macOS UI, automation scripts

---

### 13. **skill-creator** - Agent Skills Development
**Description:** Create or update AgentSkills with proper structure.

**Capabilities:**
- Skill design
- Script packaging
- Reference management
- Asset organization

**Use Cases:**
- Creating new skills
- Updating existing skills
- Skill documentation
- Skill distribution

**Integration Points:** OpenClaw skill system

---

### 14. **slack** - Slack Integration
**Description:** Control Slack from OpenClaw via the slack tool.

**Capabilities:**
- Message reactions
- Pin/unpin items
- Channel actions

**Use Cases:**
- Team communication automation
- Slack workflow automation
- Channel management

**Integration Points:** Slack API

---

### 15. **summarize** - Content Summarization
**Description:** Summarize or extract text/transcripts from various sources.

**Capabilities:**
- URL summarization
- Podcast transcription
- Local file processing
- YouTube/video transcription

**Use Cases:**
- Research summarization
- Content extraction
- Media transcription
- Document analysis

**Integration Points:** Web content, media files, documents

---

### 16. **ui-ux-pro-max** - UI/UX Design Intelligence
**Description:** Comprehensive UI/UX design tool with extensive options.

**Capabilities:**
- 67 styles
- 96 color palettes
- 57 font pairings
- 25 chart types
- 13 tech stacks (React, Next.js, Vue, Svelte, SwiftUI, React Native, Flutter, Tailwind, shadcn/ui)

**Actions:** plan, build, create, design, implement, review, fix, improve, optimize, enhance, refactor, check

**Styles:** glassmorphism, claymorphism, minimalism, brutalism, neumorphism, bento grid, dark mode, responsive, skeuomorphism, flat design

**Use Cases:**
- UI/UX design
- Component creation
- Design system implementation
- Accessibility review
- Responsive design

**Integration Points:** shadcn/ui MCP, frontend frameworks

---

### 17. **video-frames** - Video Processing
**Description:** Extract frames or short clips from videos using ffmpeg.

**Capabilities:**
- Frame extraction
- Short clip creation
- Video processing

**Use Cases:**
- Video analysis
- Thumbnail generation
- Video editing
- Frame-by-frame review

**Integration Points:** ffmpeg, video files

---

### 18. **wacli** - WhatsApp Management
**Description:** WhatsApp CLI for messaging and history management.

**Capabilities:**
- Send messages
- Search/sync history

**Use Cases:**
- WhatsApp automation
- Message history analysis
- Communication workflows

**Integration Points:** WhatsApp API

---

### 19. **weather** - Weather Information
**Description:** Get current weather and forecasts via wttr.in or Open-Meteo.

**Capabilities:**
- Current weather
- Weather forecasts
- Location-based queries

**Use Cases:**
- Weather monitoring
- Location-aware applications
- Forecast integration

**Integration Points:** wttr.in, Open-Meteo API

**Limitations:** No historical data, severe weather alerts, or detailed meteorological analysis

---

### 20. **xurl** - X (Twitter) API
**Description:** CLI for authenticated X (Twitter) API requests.

**Capabilities:**
- Post tweets, reply, quote
- Search and read posts
- Manage followers
- Send DMs
- Upload media
- Full X API v2 access

**Use Cases:**
- Twitter automation
- Social media management
- Content distribution
- Engagement tracking

**Integration Points:** X (Twitter) API v2

---

## Integration Patterns

### Common Integration Points

1. **APIs:** Most skills integrate with external APIs (GitHub, Twitter, Slack, etc.)
2. **CLIs:** Many skills wrap command-line tools
3. **Web Services:** HTTP/HTTPS endpoints
4. **Local Tools:** System utilities (ffmpeg, git, etc.)
5. **OpenClaw Core:** Direct integration with OpenClaw's tool system

### Workflow Automation

Skills can be combined for complex workflows:
- **Content Pipeline:** AI Trend Hunter Pro → ui-ux-pro-max → coding-agent → github
- **DevOps:** healthcheck → github → coding-agent → slack (notifications)
- **Communication:** summarize → himalaya → slack → imsg
- **Research:** web_search → web_fetch → summarize → gemini

---

## Recommendations for Agent-Todo Platform

### High-Value Skills for Task Management

1. **coding-agent** - Automate bug fixes and feature implementation
2. **github** - Integrate with GitHub issues for bi-directional sync
3. **slack** - Send notifications for task updates
4. **summarize** - Process task descriptions and comments

### Potential Integrations

1. **GitHub Issues Sync:** Use `gh-issues` to sync tasks with GitHub issues
2. **Automated PR Creation:** Use `coding-agent` to implement task requirements
3. **Notification System:** Use `slack` to notify teams of task status changes
4. **Documentation:** Use `summarize` to generate task summaries from long descriptions

---

## Conclusion

OpenClaw provides a rich ecosystem of 20+ skills covering:
- **Development:** coding-agent, github, gh-issues
- **Communication:** slack, imsg, wacli, himalaya
- **Content:** AI Trend Hunter Pro, summarize, gemini
- **UI/UX:** 21st-dev, ui-ux-pro-max
- **Utilities:** weather, video-frames, nano-pdf
- **System:** healthcheck, mcporter, peekaboo
- **Platform:** skill-creator

Each skill is modular and can be combined for complex automation workflows. The Agent-Todo platform can leverage these skills to enhance task management, automate workflows, and provide better integration with development tools.

---

**Next Steps:**
1. Evaluate high-priority integrations (GitHub, Slack)
2. Design workflow automation patterns
3. Implement skill-based task automation
4. Document integration best practices

---

*Generated by Agent-Todo*  
*Task ID: 2ab1df90-65e7-434b-abd1-20e2f8321e28*
