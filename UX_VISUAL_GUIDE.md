# Agent Todo Platform - UX Visual Guide

## Card Layout Examples

### Task Card Visual Hierarchy

```
┌─────────────────────────────────────────────┐ ← 4px colored border
│ ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ │ ← Critical banner (red)
│                                             │
│  Process Customer Data        🔵 In Progress │ ← Status badge
│  🔥 HIGH PRIORITY                            │ ← Priority badge
│                                             │
│  🤖 EA                                      │ ← Agent avatar (blue circle)
│  Example Agent                              │
│                                             │
│  ────────────────────────────────────────   │
│                                             │
│  Analyze and process the uploaded CSV       │
│  files to extract customer insights and...   │
│  [show more]                                │
│                                             │
│  📅 Due: Today        👤 admin@example.com  │
│                                             │
│  ▓▓▓▓▓▓░░░░░░░ 50%                         │ ← Progress bar
│                                             │
│  ────────────────────────────────────────   │
│                                             │
│  [View Details] [Edit]        💬 3  📜 5    │
└─────────────────────────────────────────────┘
```

### Agent Color Coding Examples

#### Example Agent (Blue)
```
┌─────────────────────────────────────────────┐
│ ▓▓▓▓  ← Blue border (4px)                  │
│  Setup Environment                          │
│                                             │
│  🔵 Example Agent                           │
│  [EA]  ← Blue circle with white text       │
│                                             │
│  Configure dev environment...               │
└─────────────────────────────────────────────┘

Colors:
  Primary:  #3B82F6 (Blue 500)
  Light:    #DBEAFE (Blue 100)
  Dark:     #1E40AF (Blue 800)
```

#### Data Agent (Purple)
```
┌─────────────────────────────────────────────┐
│ ▓▓▓▓  ← Purple border (4px)                │
│  Analyze Metrics                           │
│                                             │
│  🟣 Data Agent                             │
│  [DA]  ← Purple circle with white text     │
│                                             │
│  Process analytics data...                  │
└─────────────────────────────────────────────┘

Colors:
  Primary:  #A855F7 (Purple 500)
  Light:    #F3E8FF (Purple 100)
  Dark:     #6B21A8 (Purple 800)
```

#### Support Agent (Green)
```
┌─────────────────────────────────────────────┐
│ ▓▓▓▓  ← Green border (4px)                 │
│  Handle Support Ticket                     │
│                                             │
│  🟢 Support Agent                          │
│  [SA]  ← Green circle with white text      │
│                                             │
│  Respond to customer inquiry...             │
└─────────────────────────────────────────────┘

Colors:
  Primary:  #10B981 (Emerald 500)
  Light:    #D1FAE5 (Emerald 100)
  Dark:     #065F46 (Emerald 800)
```

### Status Card Variations

#### PENDING (Yellow gradient)
```
┌─────────────────────────────────────────────┐
│ ▓▓▓▓  ← Yellow border                      │
│ ═══════════════════════════════════════════ │
│                                             │
│  🟡 PENDING                                 │
│  Waiting to be started                      │
└─────────────────────────────────────────────┘

Background: linear-gradient(to right, #FFFBEB, white)
Border: #F59E0B (Amber 500)
```

#### IN PROGRESS (Blue gradient + pulse)
```
┌─────────────────────────────────────────────┐
│ ▓▓▓▓  ← Blue border                        │
│ ═══════════════════════════════════════════ │
│                                             │
│  🔵 IN PROGRESS                             │
│  [Working on it...]                         │
│                                             │
│  ▓▓▓▓▓▓░░░░░░ 50%  ← Animated progress      │
└─────────────────────────────────────────────┘

Background: linear-gradient(to right, #DBEAFE, white)
Border: #3B82F6 (Blue 500)
Animation: pulse-border (subtle)
```

#### COMPLETED (Green gradient + dimmed)
```
┌─────────────────────────────────────────────┐
│ ▓▓▓▓  ← Green border                        │
│ ═══════════════════════════════════════════ │
│                                             │
│  ✅ COMPLETED                               │
│  Finished on Jan 15, 2024                   │
└─────────────────────────────────────────────┘

Background: linear-gradient(to right, #D1FAE5, white)
Border: #10B981 (Emerald 500)
Opacity: 0.8
```

#### FAILED (Red gradient)
```
┌─────────────────────────────────────────────┐
│ ▓▓▓▓  ← Red border                          │
│ ═══════════════════════════════════════════ │
│                                             │
│  ❌ FAILED                                  │
│  Error: Connection timeout                  │
└─────────────────────────────────────────────┘

Background: linear-gradient(to right, #FEE2E2, white)
Border: #EF4444 (Red 500)
```

### Priority Indicators

#### CRITICAL
```
┌─────────────────────────────────────────────┐
│ ═══════════════════════════════════════════ │ ← Red banner
│ 🔥 CRITICAL                                 │
│ ═══════════════════════════════════════════ │
│  Task title                                 │
└─────────────────────────────────────────────┘
Banner: #EF4444 (Red 500)
Icon: 🔥
```

#### HIGH
```
┌─────────────────────────────────────────────┐
│  Task title              [⬆️ HIGH]         │
│                          ↑ Orange badge     │
└─────────────────────────────────────────────┘
Badge: #F59E0B (Amber 500)
```

#### MEDIUM
```
┌─────────────────────────────────────────────┐
│  Task title              [→ MEDIUM]        │
│                          ↑ Blue badge       │
└─────────────────────────────────────────────┘
Badge: #3B82F6 (Blue 500)
```

#### LOW
```
┌─────────────────────────────────────────────┐
│  Task title              [↓ LOW]           │
│                          ↑ Gray badge       │
└─────────────────────────────────────────────┘
Badge: #9CA3AF (Gray 400)
```

### Grid Layouts

#### Mobile (1 column)
```
┌─────────────────┐
│  Task Card 1    │
└─────────────────┘
┌─────────────────┐
│  Task Card 2    │
└─────────────────┘
┌─────────────────┐
│  Task Card 3    │
└─────────────────┘
```

#### Tablet (2 columns)
```
┌─────────────────┐ ┌─────────────────┐
│  Task Card 1    │ │  Task Card 2    │
└─────────────────┘ └─────────────────┘
┌─────────────────┐ ┌─────────────────┐
│  Task Card 3    │ │  Task Card 4    │
└─────────────────┘ └─────────────────┘
```

#### Desktop (3 columns)
```
┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│ Task Card 1 │ │ Task Card 2 │ │ Task Card 3 │
└─────────────┘ └─────────────┘ └─────────────┘
┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│ Task Card 4 │ │ Task Card 5 │ │ Task Card 6 │
└─────────────┘ └─────────────┘ └─────────────┘
```

#### Wide Desktop (4 columns)
```
┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐
│ Card 1  │ │ Card 2  │ │ Card 3  │ │ Card 4  │
└─────────┘ └─────────┘ └─────────┘ └─────────┘
┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐
│ Card 5  │ │ Card 6  │ │ Card 7  │ │ Card 8  │
└─────────┘ └─────────┘ └─────────┘ └─────────┘
```

### Hover Effects

#### Normal State
```
┌─────────────────────────────────┐
│  Task Title                     │
│  box-shadow: 0 1px 3px rgba(0)  │
│  transform: translateY(0)       │
└─────────────────────────────────┘
```

#### Hover State
```
┌─────────────────────────────────┐
│  Task Title                     │
│  box-shadow: 0 12px 24px rgba(0)│ ← Lift effect
│  transform: translateY(-4px)    │ ← Move up
│  transition: all 0.3s ease     │ ← Smooth
└─────────────────────────────────┘
```

### Agent Avatar Sizes

#### Large (Card Header)
```
┌──────────┐
│          │
│    EA    │  40px diameter
│          │
└──────────┘
```

#### Medium (Inline)
```
┌──────┐
│  EA  │  32px diameter
└──────┘
```

#### Small (Compact)
```
┌────┐
│ E  │  24px diameter
└────┘
```

### Filter Bar Design

```
┌─────────────────────────────────────────────────────────────┐
│ 🔍 [Search tasks...]                        [⚙️ Filters]   │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  Status:     [All ▼] [Pending] [In Progress] [Done]       │
│                                                              │
│  Priority:  [All ▼] [Critical] [High] [Medium] [Low]       │
│                                                              │
│  Agent:      [All ▼]                                         │
│              ┌─────┐ [🤖 Example Agent] (5 tasks)          │
│              │  EA  │  ← Blue avatar                        │
│              └─────┘                                       │
│              ┌─────┐ [🤖 Data Agent] (3 tasks)             │
│              │  DA  │  ← Purple avatar                       │
│              └─────┘                                       │
│                                                              │
│  Date:       [Any time ▼] [Today] [This week] [This month] │
│                                                              │
│  [Clear All Filters]                            [Apply]    │
└─────────────────────────────────────────────────────────────┘
```

### Kanban Board View (Optional)

```
┌──────────────┐  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│  PENDING     │  │ IN PROGRESS  │  │  COMPLETED   │  │   FAILED     │
│    (3)       │  │     (2)      │  │     (5)      │  │     (1)      │
├──────────────┤  ├──────────────┤  ├──────────────┤  ├──────────────┤
│              │  │              │  │              │  │              │
│ ┌──────────┐ │  │ ┌──────────┐ │  │ ┌──────────┐ │  │ ┌──────────┐ │
│ │ Task 1   │ │  │ │ Task 4   │ │  │ │ Task 2   │ │  │ │ Task 7   │ │
│ └──────────┘ │  │ └──────────┘ │  │ └──────────┘ │  │ └──────────┘ │
│              │  │              │  │              │  │              │
│ ┌──────────┐ │  │ ┌──────────┐ │  │ ┌──────────┐ │  │              │
│ │ Task 3   │ │  │ │ Task 5   │ │  │ │ Task 6   │ │  │              │
│ └──────────┘ │  │ └──────────┘ │  │ └──────────┘ │  │              │
│              │  │              │  │              │  │              │
│ ┌──────────┐ │  │              │  │ ┌──────────┐ │  │              │
│ │ Task 8   │ │  │              │  │ │ Task 9   │ │  │              │
│ └──────────┘ │  │              │  │ └──────────┘ │  │              │
│              │  │              │  │              │  │              │
│ ┌──────────┐ │  │              │  │ ┌──────────┐ │  │              │
│ │ Task 10  │ │  │              │  │ │ Task 11  │ │  │              │
│ └──────────┘ │  │              │  │ └──────────┘ │  │              │
│              │  │              │  │              │  │              │
│ [+ Add Task] │  │ [+ Add Task] │  │ [+ Add Task] │  │ [+ Add Task] │
└──────────────┘  └──────────────┘  └──────────────┘  └──────────────┘
```

---

## Color Swatches

### Agent Colors (Auto-generated)
```
1.  Blue:    #3B82F6 (Hue 220)
2.  Green:   #10B981 (Hue 160)
3.  Orange:  #F59E0B (Hue 30)
4.  Pink:    #EC4899 (Hue 330)
5.  Teal:    #14B8A6 (Hue 180)
6.  Red:     #EF4444 (Hue 0)
7.  Purple:  #8B5CF6 (Hue 270)
8.  Yellow:  #EAB308 (Hue 45)
9.  Sky:     #0EA5E9 (Hue 200)
10. Violet:  #7C3AED (Hue 280)
```

### Status Colors
```
Pending:     #F59E0B (Amber 500)
In Progress: #3B82F6 (Blue 500)
Completed:   #10B981 (Emerald 500)
Failed:      #EF4444 (Red 500)
```

### Priority Colors
```
Critical: #EF4444 (Red 500)
High:     #F59E0B (Amber 500)
Medium:   #3B82F6 (Blue 500)
Low:      #9CA3AF (Gray 400)
```

---

## Animation Timings

### Card Entry
- Duration: 0.3s
- Easing: ease-out
- Stagger: 50ms between cards

### Hover
- Duration: 0.3s
- Easing: ease-in-out
- Delay: 0s

### Status Change
- Duration: 0.6s
- Type: pulse
- Iterations: 1

### Filter Transition
- Duration: 0.3s
- Easing: ease-in-out

---

## Accessibility

### Color Contrast Ratios
- ✅ All text ≥ 4.5:1 (WCAG AA)
- ✅ Large text ≥ 3:1 (WCAG AA)
- ✅ UI components ≥ 3:1 (WCAG AA)

### Keyboard Navigation
- Tab through cards in order
- Enter/Space to expand card
- Arrow keys in grid view
- Escape to close details

### Screen Reader Support
- Card: "article" role
- Status: "status" landmark
- Agent: "img" role with alt text
- Actions: "button" roles

---

This visual guide provides a complete reference for implementing the card-based UX with agent colors! 🎨
