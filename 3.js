// File: 3.js
// Basic JavaScript file with module structure

console.log("3.js file is working!");

// Example function for demonstration
function greetUser(name = "World") {
    return `Hello, ${name}!`;
}

// Basic module exports
module.exports = {
    greetUser
};

// Demonstrate functionality
console.log(greetUser());
console.log(greetUser("Developer"));
