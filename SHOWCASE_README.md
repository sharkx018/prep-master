# 🚀 PrepMaster Pro - Interview Preparation Platform

<div align="center">
  <img src="https://img.shields.io/badge/Version-1.0-blue.svg" alt="Version 1.0">
  <img src="https://img.shields.io/badge/Built%20with-React%20%26%20Go-brightgreen.svg" alt="Built with React & Go">
  <img src="https://img.shields.io/badge/Database-PostgreSQL-blue.svg" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/Styled%20with-Tailwind%20CSS-06B6D4.svg" alt="Tailwind CSS">
</div>

<div align="center">
  <h3>🎯 Master Your Technical Interviews with Style</h3>
  <p><strong>Engineered by Mukul Verma</strong></p>
</div>

---

## ✨ Overview

**PrepMaster Pro** is a modern, full-stack interview preparation platform designed to help software engineers systematically prepare for technical interviews. Track your progress across Data Structures & Algorithms, Low-Level Design, and High-Level Design topics with a beautiful, intuitive interface.

## 🎨 Features

### 📊 **Intelligent Dashboard**
- Real-time progress tracking with visual indicators
- Completion cycles counter
- Motivational messages based on your progress
- Quick action cards for seamless navigation

### 📚 **Smart Study Mode**
- Random problem selection algorithm
- Progress bar visible during study sessions
- Mark items complete with instant feedback
- Skip functionality for flexible learning

### 🗂️ **Comprehensive Item Management**
- Add problems/articles with categorization
- Filter by category, subcategory, and status
- Bulk operations support
- Beautiful card-based interface

### 📈 **Advanced Analytics**
- Interactive charts (Pie & Bar charts)
- Category-wise breakdown
- Subcategory progress grids
- Completion cycle tracking

### 🎯 **Categories Supported**
- **DSA**: Arrays, Strings, Trees, Graphs, Dynamic Programming, and more
- **LLD**: Design Patterns, SOLID Principles, OOP Concepts
- **HLD**: Distributed Systems, Microservices, System Design

## 🛠️ Tech Stack

### Frontend
- **React 18** with TypeScript
- **Tailwind CSS** for modern styling
- **Recharts** for data visualization
- **React Router** for navigation
- **Axios** for API communication
- **Lucide React** for beautiful icons

### Backend
- **Go (Golang)** with Gin framework
- **PostgreSQL** database
- **Clean Architecture** pattern
- **RESTful API** design
- **Docker** for containerization

## 🚀 Quick Start

### Prerequisites
- Node.js 16+
- Go 1.19+
- PostgreSQL 14+
- Docker (optional)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/mukulverma/prepmaster-pro.git
   cd prepmaster-pro
   ```

2. **Start the database**
   ```bash
   make db-only
   ```

3. **Run the backend**
   ```bash
   go run cmd/server/main.go
   ```

4. **Start the frontend**
   ```bash
   cd frontend
   npm install
   npm start
   ```

5. **Access the app**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

## 📸 Screenshots

### Dashboard
<div align="center">
  <img src="screenshots/dashboard.png" alt="Dashboard" width="600">
</div>

### Study Mode
<div align="center">
  <img src="screenshots/study.png" alt="Study Mode" width="600">
</div>

### Analytics
<div align="center">
  <img src="screenshots/analytics.png" alt="Analytics" width="600">
</div>

## 🌟 Key Highlights

- **🎨 Modern UI/UX**: Gradient designs, smooth animations, and responsive layouts
- **⚡ High Performance**: Optimized queries and efficient state management
- **📱 Fully Responsive**: Works seamlessly on desktop and mobile devices
- **🔒 Clean Code**: Following best practices and design patterns
- **🚀 Production Ready**: Docker support and environment configuration

## 📝 API Documentation

### Items Endpoints
- `POST /api/v1/items` - Create new item
- `GET /api/v1/items` - List items with filters
- `GET /api/v1/items/next` - Get random pending item
- `PUT /api/v1/items/:id/complete` - Mark item as complete
- `DELETE /api/v1/items/:id` - Delete item

### Statistics Endpoints
- `GET /api/v1/stats` - Overall statistics
- `GET /api/v1/stats/detailed` - Detailed breakdown
- `GET /api/v1/stats/category/:category` - Category statistics

## 🤝 Contributing

While this is a personal project by Mukul Verma, suggestions and feedback are welcome! Feel free to open issues or submit pull requests.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 👨‍💻 Author

**Mukul Verma**
- GitHub: [@mukulverma](https://github.com/mukulverma)
- LinkedIn: [Mukul Verma](https://linkedin.com/in/mukulverma)

---

<div align="center">
  <p>Built with ❤️ by Mukul Verma</p>
  <p>⭐ Star this repository if you find it helpful!</p>
</div> 