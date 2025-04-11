/**
 * A simple JavaScript file with example code
 */

// Variable declaration
const username = "User";

/**
 * Returns a greeting message with the given name
 * @param {string} name - The name to greet
 * @return {string} The greeting message
 */
function getGreeting(name) {
  return `Hello, ${name}! Welcome to the application.`;
}

// Demonstrate function call with console.log
console.log(getGreeting(username));