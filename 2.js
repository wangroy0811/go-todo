/**
 * Simple calculator JavaScript example
 */
222222
333333333333
// Variable declarations
const numberA = 10;
const numberB = 5;

/**
 * Calculator object with basic arithmetic operations
 */
const calculator = {
  /**
   * Adds two numbers
   * @param {number} a - First number
   * @param {number} b - Second number
   * @return {number} Sum of a and b
   */
  add: function(a, b) {
    return a + b;
  },
  
  /**
   * Subtracts second number from first number
   * @param {number} a - First number
   * @param {number} b - Second number to subtract
   * @return {number} Difference of a and b
   */
  subtract: function(a, b) {
    return a - b;
  }
};

/**
 * Formats calculation result with operation description
 * @param {number} a - First operand
 * @param {number} b - Second operand
 * @param {string} operation - Name of operation performed
 * @param {number} result - Result of calculation
 * @return {string} Formatted result string
 */
function formatResult(a, b, operation, result) {
  return `${a} ${operation} ${b} = ${result}`;
}

// Demonstrate calculator functionality
const addResult = calculator.add(numberA, numberB);
const subtractResult = calculator.subtract(numberA, numberB);

console.log(formatResult(numberA, numberB, '+', addResult));
console.log(formatResult(numberA, numberB, '-', subtractResult));