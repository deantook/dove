# design.md

## Design Overview

产品整体风格属于：

- 极简主义（Minimalism）
- 大留白（Whitespace-first）
- 内容驱动（Content-first）
- 柔和社交产品风格
- iOS Human Interface Guidelines 倾向

关键词：

`Clean` · `Warm` · `Friendly` · `Editorial`

---

# 1. Color System

## Primary Color

Hinge 的核心品牌色为紫色系按钮色。

| Token | Color |
|---------|---------|
| Primary | #7C3AED |
| Primary Hover | #6D28D9 |
| Primary Light | #EDE9FE |

主要用于：

- CTA Button
- 底部导航激活态
- Tag选中状态
- 重点操作入口

---

## Secondary Colors

### Purple Tint

```css
#A78BFA
```

用于：

- 次级强调
- Icon高亮
- 筛选器状态

---

### Soft Pink

```css
#F5E6EC
```

用于：

- 插画背景
- Empty State
- Onboarding

---

## Background Colors

### Main Background

```css
#FFFFFF
```

占比极高。

---

### Secondary Background

```css
#F7F7F7
```

用于：

- Card容器
- 输入区域
- 底部弹层

---

### Divider Background

```css
#EFEFEF
```

用于：

- 分割线
- 卡片边界

---

## Text Colors

### Primary Text

```css
#111111
```

主要标题。

---

### Secondary Text

```css
#666666
```

说明文案。

---

### Tertiary Text

```css
#9B9B9B
```

Placeholder、辅助信息。

---

### Disabled Text

```css
#C7C7C7
```

不可用状态。

---

# 2. Typography

整体接近 iOS SF Pro Display / SF Pro Text。

---

## Page Title

```css
font-size: 32px;
font-weight: 700;
line-height: 40px;
```

示例：

- What's your name?
- You're one of a kind.

---

## Section Title

```css
font-size: 24px;
font-weight: 600;
line-height: 32px;
```

---

## Card Headline

```css
font-size: 18px;
font-weight: 600;
line-height: 26px;
```

用于：

- Profile Prompt
- 卡片标题

---

## Body Text

```css
font-size: 16px;
font-weight: 400;
line-height: 24px;
```

这是最常见字号。

---

## Caption

```css
font-size: 14px;
font-weight: 400;
line-height: 20px;
```

用于：

- 标签
- 辅助说明
- 设置页面

---

## Small Meta Text

```css
font-size: 12px;
font-weight: 400;
line-height: 16px;
```

用于：

- 状态提示
- Secondary Info

---

# 3. Spacing System

整体遵循 4pt Grid，但核心节奏更接近 8pt Grid。

---

## Base Scale

```css
4px
8px
12px
16px
24px
32px
40px
48px
64px
```

---

## Page Padding

左右边距：

```css
20px
```

部分页面：

```css
24px
```

---

## Card Padding

```css
16px
```

较大内容区：

```css
20px
24px
```

---

## Component Spacing

输入框与标题：

```css
12px
```

按钮与内容：

```css
16px
```

模块与模块：

```css
24px
```

大区块：

```css
32px
```

---

## Visual Breathing Room

Hinge 最大特征之一：

大量留白。

典型比例：

```text
内容高度：留白高度 ≈ 1 : 0.6
```

即使只有一个输入框，也会保留大量空白区域。

---

# 4. Component Style

## Primary Button

### Size

```css
height: 48px;
```

---

### Radius

```css
border-radius: 24px;
```

胶囊按钮。

---

### Style

```css
background: #7C3AED;
color: white;
font-weight: 600;
```

---

## Secondary Button

```css
background: transparent;
border: none;
color: #111111;
```

文字按钮风格。

---

## Text Input

### Height

```css
48px
```

---

### Border

```css
1px solid #E5E5E5;
```

---

### Radius

```css
12px
```

---

## Card

### Radius

```css
16px
```

照片卡片：

```css
16px ~ 20px
```

---

### Border

```css
1px solid #EFEFEF;
```

通常非常浅。

---

### Shadow

```css
0 2px 8px rgba(0,0,0,0.05)
```

极弱阴影。

很多页面甚至无阴影。

---

## Bottom Sheet

### Radius

```css
24px 24px 0 0
```

---

### Background

```css
#FFFFFF
```

---

### Shadow

```css
0 -4px 20px rgba(0,0,0,0.08)
```

---

## Navigation Bar

高度：

```css
56px
```

---

### Active State

```css
color: #7C3AED;
```

---

### Inactive State

```css
color: #999999;
```

---

### Indicator

无明显背景块。

仅使用：

- Icon颜色变化
- 细线高亮

---

# 5. Border Radius System

统一偏圆润。

```css
8px
12px
16px
20px
24px
```

使用频率：

| Radius | Usage |
|----------|----------|
| 8px | 小标签 |
| 12px | 输入框 |
| 16px | 卡片 |
| 20px | 图片容器 |
| 24px | CTA按钮 |

---

# 6. Elevation System

非常克制。

### Level 0

```css
box-shadow: none;
```

### Level 1

```css
0 2px 8px rgba(0,0,0,.05)
```

### Level 2

```css
0 8px 24px rgba(0,0,0,.08)
```

主要用于弹层。

---

# 7. Overall Design Rules

1. 白色背景占比 > 75%
2. 品牌紫只出现在关键 CTA
3. 字体层级少，不超过 5 个等级
4. 卡片圆角大于 16px
5. 阴影极弱
6. 大量留白优先于信息密度
7. 内容（照片/Prompt）优先于装饰元素
8. 几乎不使用纯黑，使用 #111111
9. 使用浅灰分隔，而非重边框
10. CTA 永远最显眼