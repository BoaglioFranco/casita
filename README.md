The goal of this project is to create a quick app that lets my wife and I share a groceries list and get push notifications each time one adds or removes an item from the list. The push notifications are very important.
We have 2 iphones, i don't intend to share this app with anyone and i don't have a developer account. So i came to the conclusion that the best way to implement this is with a PWA. 
The backend is a small Go server backed by PostgreSQL, and the UI is built with React. The PWA push notification flow is the core feature.

Both services apply pending database migrations on startup. See
[database migrations](docs/database-migrations.md) for setup and the one-time
baseline required for databases created before migration tracking was introduced.
