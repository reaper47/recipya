---
sidebar_position: 2
---

# Configuration File

The [configuration file](https://github.com/reaper47/recipya/blob/main/deploy/config.example.json)
sets important variables for the application. Let's go over each of them.

  | **Variable**         | **Description**                                                                                                       |
  |----------------------|-----------------------------------------------------------------------------------------------------------------------|
  | email.from           | The administrator's email address. It is usually the email address of your [SendGrid](https://sendgrid.com/) account. |
  | email.sendGridAPIKey | Your [SendGrid](https://sendgrid.com/) API key. The free tier should be sufficient for your needs.                    |
  | isProduction         | Whether the app is in production. Its value is either **true** or **false**.                                          |
  | port                 | The port the app will be served through if localhost.                                                                 |
  | url                  | The website the app is served on. This URL will serve as the base link in the emails.                                 |
