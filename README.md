# Tukang Catat Uang Bot

This repository contains a Telegram bot built with Golang that interacts with users to manage transactions, categories, and accounts via a simple and intuitive interface.

## Features

- **Start Command**: Initiate the bot and start interacting.
- **Income/Expense Management**: Easily record income or expenses.
- **Transfer Management**: Handle account transfers.
- **Reporting**: Generate transaction reports.

## Getting Started

To get a local copy up and running, follow these steps.

### Prerequisites

- Go 1.20+ installed on your machine
- A PostgreSQL database instance
- A Telegram bot token (you can get one from [BotFather](https://core.telegram.org/bots#botfather))

### Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/yourusername/telegram-bot.git
    cd telegram-bot
    ```

2. **Install the necessary dependencies:**

    ```bash
    go mod tidy
    ```

3. **Set up your environment variables:**

   Create a `.env` file in the root directory and add your configuration values:

    ```dotenv
    SERVER_HOST=0.0.0.0
    SERVER_REST_PORT=8080
    TELEGRAM_BOT_TOKEN=**************************
    TELEGRAM_USERNAME_ALLOWED=***,***,***
    TGBOT_DB_POSTGRESQL_URL=**************************
    ```

4. **Run the bot:**

    ```bash
    go run main.go
    ```

### Accessing the Bot

- The bot will be available on Telegram at the bot token URL you configured.
- The HTTP server will be accessible at `http://localhost:8080/`.

### Contribution

We welcome contributions! Please follow the guidelines below:

1. **Fork the Project**

   Click the Fork button at the top right of this repository to fork your own version.

2. **Create a Feature Branch**

   Create a branch for your feature or bugfix:

    ```bash
    git checkout -b feature/your-feature-name
    ```

3. **Commit Your Changes**

   Make your changes and commit them with a descriptive message:

    ```bash
    git commit -m 'Add some feature'
    ```

4. **Push to the Branch**

   Push your changes to your forked repository:

    ```bash
    git push origin feature/your-feature-name
    ```

### License

This project is licensed under the MIT License - see the `LICENSE` file for details.

### Acknowledgments

Thanks to the creators of the libraries and tools used in this project.
Special thanks to BotFather for providing an easy way to create and manage Telegram bots.