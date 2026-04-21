package telegram

const startMsg = `Welcome to the Finance Tracker Bot! 💰

✅ Create or edit categories  
✅ Track balance  
✅ Update categories automatically  

Choose an option below: ⬇️`

const listCategoriesBalanceMsg = "🐶 Name:                    💶 Balance:        "
const listCategoriesPermanentBalanceMsg = "🐶 Name:                     📛 Limit:            "
const listCategoriesDeleteMsg = "🐶 Name:"
const deleteCategoryMsg = "⚠️ Are you sure you want to delete this category? This action cannot be undone! ❌"
const updateBalancesMsg = "⚠️ Are you sure you want to update balances? This action cannot be undone! ❌"
const deleteCategorySuccessMsg = "The category '%s' has been successfully deleted. ✅"
const defaultCallbackMsg = "Oops! 😕 Something went wrong while processing your request. Please try again later."
const defaultCommandMsg = "🚀 Oops! This command doesn’t exist. Try again! 😉"
const createCategoryMsg = "📂 Enter the category name 👇👇👇"
const renameCategoryMsg = "🔄 Enter a new category name 👇👇👇"
const enterReceiptMsg = "💰 Enter the receipt amount 👇👇👇"
const enterLimitMsg = "💰 Enter the limit amount 👇👇👇"
const enterAllowanceBalanceMsg = "💸 Please enter the amount you plan to spend this month 👇👇👇"
const incorrectNumberBalanceMsg = "⚠️ Please enter a numerical value for the balance."
const updateBalanceSuccessMsg = "🎉 Balance for category '%s' updated successfully! 💸"
const updateLimitSuccessMsg = "🎉 Limit for category '%s' updated successfully! 💸"
const categoryCreatedMsg = "🎉 Category '%s' created successfully!"
const categoryRenamedMsg = "🔄 Category '%s' has been renamed to '%s' successfully! 🎉"
const unknownErrorMsg = "❗️ An unknown error occurred. Please try again later."
const updateBalancesSuccessMsg = "🎉 Balances for all categories have been updated successfully! 💸"
const categoryNotFoundMsg = "Oops! 😢 It looks like this category no longer exists...You can create a new one and keep going! 🚀"
const categorySelectedMsg = "✅ Category '%s' selected! Now, choose an action below: ⬇️"
