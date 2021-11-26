package telegramBot

const startMessage = `
Welcome to kraken futures trading bot! 📈

This bot supports trading on simple indicator calling "stop loss & take profit"

Here is bot commands:
	💁 /help - list descriptions for commands
`

const helpMessage = `
Commands 📟:
	🔵 /help - list descriptions for commands
	🔵 /sign_up - registrate you in trading bot system
	🔵 /exit_from_sign_up - stop getting input data to register you in the bot
	🔵 /sign_in - login you in trading bot system and allow to trade on kraken futures
	🔵 /exit_from_sign_in - stop getting input data to login you in the bot
	🔵 /send_order - allow to send market order with symbol, side and amount arguments to kraken futures
	🔵 /exit_from_send_order - stop getting input data to send order to kraken futures
	🔵 /logout - logout you from trading bot system on every telegram device associated with your username
`

const invalidCommandMessage = `
⛔ No such command
`

const signUpMessage = `
🔳 Enter message in format:

Name
Username
Password
Your public api key
Your private api key

🔳 Example:

Ivan ivan password key key
`

const signUpErrMessage = `
⛔ Unable to continue further execution of sign up due to
`

const signUpSuccessMessage = `
✅ User successfully registered!
`

const signInMessage = `
🔳 Enter message in format:

Username
Password

🔳 Example:

ivan password
`

const signInErrMessage = `
⛔ Unable to continue further execution of sign in due to
`

const signInSuccessMessgae = `
✅ User successfully logged in!
`

const logoutErrMessage = `
⛔ Unable to continue further execution of logout due to
`

const logoutSuccessMessgae = `
✅ User successfully logged out!
`

const sendOrderMessage = `
🔳 Enter message in format:

Symbol (one of symbols on kraken futures)
Side   (buy or sell)      
Size   (integer up to 25000)      

🔳 Example:

PI_XBTUSD buy 10000
`

const sendOrderErrMessage = `
⛔ Unable to continue further execution of send order due to
`

const sendOrderSuccessMessage = `
✅ Successfully send order!
`
