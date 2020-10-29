/*

	This code filters existing addresses. It keeps only addresses which made swaps on uniswap.

	1. Loads addresses from a csv file.
	2. Loads transactions from etherscan.
	3. Checks if specific transaction interacted with (send money to) uniswap contract.
	4. get transaction from uniswap.
	5. if it has swaps, save address.
*/
package main

func main() {}
