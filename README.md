# 🍱 Huly.io Plugin Challenge: LunchRush

Welcome to the **LunchRush** challenge — your mission is to build a collaborative Huly.io plugin that helps employees coordinate their daily lunch orders in a fun and efficient way.

---

## 🧠 The Idea

Lunchtime coordination is always messy: scattered messages, missed orders, and last-minute chaos. **LunchRush** solves this by offering a central place — inside Huly.io — where team members can:

- Propose or vote on restaurants and dishes
- Join the group lunch order
- See what others are getting
- Nominate someone to place the order
- Lock the order at a set time and notify everyone

This should be a collaborative experience. Think "Trello meets lunch."

---

## 🧱 What You'll Build

### 🧩 A Huly Plugin (Frontend)
- Display the current day’s lunch session
- Allow users to join, select meals, and interact with others
- Show a live view of participants and their choices
- Lock the session and display the final summary

### 🛠 A Go Microservice (Backend)
- Use **[Dapr](https://dapr.io/)** building blocks for:
  - **Pub/Sub** for real-time updates across users
  - **State Store** for shared session data (e.g. Redis)
  - Optional: **Bindings** or **Secrets** to simulate 3rd party APIs

---

## 🚀 What We’re Looking For

This challenge is designed to evaluate your ability to:

- 🧠 Model collaborative workflows
- 👩‍💻 Write clean, idiomatic Go code
- ⚙️ Use Dapr to manage distributed state and pub/sub
- 🎨 Create a clear, user-friendly interface inside Huly.io

---

## 🧪 Bonus Ideas

You're welcome (but not required) to go further:

- "Reorder last week’s lunch"
- Anonymous voting or reactions
- Scheduled daily reminders
- Light gamification: who orders most often? Who’s always late?

---

## 📝 Submission Guidelines

1. Fork this repo
2. Implement your solution in:
   - `plugin/` for the Huly.io plugin
   - `microservice/` for your Go+Dapr backend
3. Include a `README.md` with:
   - Setup instructions
   - Anything you'd like us to know
4. After 1 week, we'll check the forks

---

## 🕒 Time Limit

You have **1 week** from when you accept the challenge. Don't worry about polish — we value thoughtfulness, clarity, and how you approach collaboration.

## 🧠 Ask AI for Help

Use whatever AI tool is in your toolbelt. If you don’t have any… I’ve got bad news for you 🙂

---

Happy coding and buon appetito! 🍽️
