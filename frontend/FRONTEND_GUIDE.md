# Interview Prep App - Frontend Guide

## 🚀 Getting Started

Your Interview Prep App is now running! Access it at: **http://localhost:3000**

## 📱 App Features

### 1. **Dashboard** (Home Page)
- View overall statistics
- See your progress percentage
- Quick access to key actions
- Total items, completed items, pending items, and completion cycles

### 2. **Study Mode**
- Click "Start Studying" from the dashboard
- Get random pending items one at a time
- Options:
  - **Open Link** - Opens the problem/article in a new tab
  - **Skip** - Get another random item
  - **Mark as Complete** - Mark the item as done and get the next one
- Shows a congratulations message when all items are completed

### 3. **Items Management**
- View all your items in a list
- Filter by:
  - Category (DSA, LLD, HLD)
  - Subcategory (e.g., arrays, strings)
  - Status (Pending, Completed)
- Actions for each item:
  - Open link
  - Mark as complete (if pending)
  - Delete item

### 4. **Add New Item**
- Add problems or articles to your study list
- Select category (DSA/LLD/HLD)
- Subcategory automatically updates based on category
- Add title and link
- Saves and redirects to items list

### 5. **Statistics**
- Visual charts showing your progress
- Pie chart for overall completion
- Bar chart for category breakdown
- Detailed subcategory progress grids
- Track completion cycles

## 🎯 Quick Start Guide

### Adding Your First Items

1. Click **"Add Item"** in the sidebar
2. Enter details:
   - **Title**: e.g., "Two Sum Problem"
   - **Link**: e.g., "https://leetcode.com/problems/two-sum/"
   - **Category**: Choose DSA, LLD, or HLD
   - **Subcategory**: Select from dropdown (e.g., "arrays")
3. Click **"Save Item"**

### Starting a Study Session

1. Go to **"Study"** in the sidebar
2. Click **"Get Next Item"**
3. Open the link in a new tab
4. Study/solve the problem
5. Click **"Mark as Complete"** when done
6. The next random item will appear automatically

### Tracking Progress

1. Visit **"Dashboard"** to see overall progress
2. Go to **"Statistics"** for detailed insights
3. Check subcategory progress to identify weak areas

## 🎨 UI Features

- **Responsive Design**: Works on desktop and mobile
- **Color Coding**:
  - Blue: DSA items
  - Green: LLD items
  - Purple: HLD items
  - Yellow: Pending status
  - Green: Completed status
- **Interactive Charts**: Hover over charts for details
- **Real-time Updates**: Changes reflect immediately

## 💡 Tips

1. **Bulk Add Items**: Add all your study materials at once, then use Study mode daily
2. **Focus on Categories**: Use filters to focus on one category at a time
3. **Track Subcategories**: Check statistics to see which subcategories need more work
4. **Completion Cycles**: Try to complete all items multiple times for better retention

## 🛠️ Troubleshooting

### Frontend Not Loading?
- Make sure the backend is running: `go run cmd/server/main.go`
- Check if PostgreSQL is running: `docker ps`
- Verify ports 3000 and 8080 are not in use

### API Errors?
- Check the browser console for errors
- Ensure CORS is enabled (already configured)
- Verify the backend is accessible at http://localhost:8080

### Styling Issues?
- Hard refresh the page: `Cmd+Shift+R` (Mac) or `Ctrl+Shift+R` (Windows/Linux)
- Clear browser cache if needed

## 🚦 Status Indicators

- **Loading States**: Spinning loader icons
- **Error Messages**: Red alert boxes with error details
- **Success Actions**: Immediate UI updates
- **Empty States**: Helpful messages when no data

Enjoy your interview preparation journey! 🎉 