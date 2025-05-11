#include <stdio.h>
#include <stdbool.h>
#include <stdlib.h>
#include <time.h>
#include <string.h>


// This method checks if the board is valid across, down, and within the 3x3 grid. 
bool is_valid(int board[9][9], int row, int col, int num) {
    // Check if we find the same num in the similar row
    for (int x = 0; x < 9; x++) {
        if (board[row][x] == num) {
            return false;
        }
    }

    // Check if we find the same num in the similar column
    for (int x = 0; x < 9; x++) {
        if (board[x][col] == num) {
            return false;
        }
    }

    // Check if we find the same num in the particular 3*3 matrix
    int start_row = row - row % 3;
    int start_col = col - col % 3;
    for (int i = 0; i < 3; i++) {
        for (int j = 0; j < 3; j++) {
            if (board[i + start_row][j + start_col] == num) {
                return false;
            }
        }
    }
    return true;
}

// A simple algorithm to find some easier values soduku values that can be derived from the puzzle
void fill_naked_singles(int board[9][9]) {
    for (int i = 0; i < 9; i++) {
        for (int j = 0; j < 9; j++) {
            if (board[i][j] == 0) {
                int possible_values[9] = {0};
                int count = 0;
                for (int num = 1; num <= 9; num++) {
                    if (is_valid(board, i, j, num)) {
                        possible_values[count] = num;
                        count++;
                    }
                }
                if (count == 1) {
                    board[i][j] = possible_values[0];
                }
            }
        }
    }
}

// A simple algorithm to find some easier values soduku values that can be derived from the puzzle
void fill_unique_candidates(int board[9][9]) {
    for (int i = 0; i < 9; i++) {
        for (int j = 0; j < 9; j++) {
            if (board[i][j] == 0) {
                for (int num = 1; num <= 9; num++) {
                    if (is_valid(board, i, j, num)) {
                        // Check rows
                        int row_count = 0;
                        for (int k = 0; k < 9; k++) {
                            if (is_valid(board, i, k, num)) {
                                row_count++;
                            }
                        }
                        if (row_count == 1) {
                            board[i][j] = num;
                            break;
                        }

                        // Check columns
                        int col_count = 0;
                        for (int k = 0; k < 9; k++) {
                            if (is_valid(board, k, j, num)) {
                                col_count++;
                            }
                        }
                        if (col_count == 1) {
                            board[i][j] = num;
                            break;
                        }

                        // Check 3x3 boxes
                        int start_row = i - i % 3;
                        int start_col = j - j % 3;
                        int box_count = 0;
                        for (int k = 0; k < 3; k++) {
                            for (int l = 0; l < 3; l++) {
                                if (is_valid(board, start_row + k, start_col + l, num)) {
                                    box_count++;
                                }
                            }
                        }
                        if (box_count == 1) {
                            board[i][j] = num;
                            break;
                        }
                    }
                }
            }
        }
    }
}

// solve the puzzle using a recursive back trace  
bool solve_sudoku_recursive(int board[9][9]) {
    for (int i = 0; i < 9; i++) {
        for (int j = 0; j < 9; j++) {
            if (board[i][j] == 0) {
                for (int num = 1; num <= 9; num++) {
                    if (is_valid(board, i, j, num)) {
                        board[i][j] = num;
                        if (solve_sudoku_recursive(board)) {
                            return true;
                        }
                        board[i][j] = 0;
                    }
                }
                return false;
            }
        }
    }

    return true;
}

// solve the puzzle
bool solve_sudoku(int board[9][9]) {
    fill_naked_singles(board);
    fill_unique_candidates(board);
    return solve_sudoku_recursive(board);

}

// count how many valid solutions exist for this puzzle
int count_solutions_recursive(int board[9][9]) {
    int count = 0;
    for (int i = 0; i < 9; i++) {
        for (int j = 0; j < 9; j++) {
            if (board[i][j] == 0) {
                for (int num = 1; num <= 9; num++) {
                    if (is_valid(board, i, j, num)) {
                        board[i][j] = num;
                        count += count_solutions_recursive(board);
                        board[i][j] = 0;
                    }
                }
                return count;
            }
        }
    }
    return 1; // if we reach this point, it means we've found a valid solution
}

void print_board(int board[9][9]) {
  printf("\n");
    for (int i = 0; i < 9; i++) {
        for (int j = 0; j < 9; j++) {
            printf("%d ", board[i][j]);
        }
        printf("\n");
    }
}

// receives a multi solution puzzle and one of the solutions and uses the solution to add values to the puzzle until it has one solution
bool correct_puzzle(int puzzle[9][9], int solution[9][9]) {
    int tries = 0;
    int sol_count = 99;
    // Check how many solution it has to begin with
    sol_count = count_solutions_recursive(puzzle);
    if(sol_count == 1){
      return true;
    }
    if(sol_count < 1){
      return false;
    }
    // try up to 1 million times to randomly select a value from the puzzle that is empty and copy the value from the solution.
    while (tries < 9999999) {
      int rand_x = rand() % (9 - 0 + 1) + 0;
      int rand_y = rand() % (9 - 0 + 1) + 0;
      if (puzzle[rand_x][rand_y] == 0){
        puzzle[rand_x][rand_y] = solution[rand_x][rand_y];
        sol_count = count_solutions_recursive(puzzle);
        if(sol_count == 1){
          return true;
        }
      }
      tries++;
    }
    return false;
}

// Returns a single solution puzzle given a multi solution puzzle
int (*get_puzzle(int puzzle[9][9]))[9] {
    static int corrected_board[9][9];
    int solved_board[9][9];

    memcpy(solved_board, puzzle, sizeof(int) * 9 * 9);  // Copy the initial puzzle to solved_board
    int solutions_count = count_solutions_recursive(puzzle);
    // If there is no solution then return an empty board
    if (solutions_count < 1 || !solve_sudoku(solved_board)) {    
        memset(corrected_board, 0, sizeof(corrected_board));  // Set to all zeros
        return corrected_board;
    }

    // copy the original board to the corrected board as a starting point
    memcpy(corrected_board, puzzle, sizeof(int) * 9 * 9);

    if (solutions_count == 1) {
        return corrected_board;
    }

    // Correct the puzzle and return the corrected version
    correct_puzzle(corrected_board, solved_board);
    return corrected_board;
}

int main() {

  // TEST 1, check that a multi solution puzzle, returns a single solution puzzle
    int input_1[9][9] = {
        {0, 6, 0, 7, 0, 5, 1, 0, 4},
        {4, 5, 0, 0, 0, 0, 0, 9, 0},
        {3, 0, 0, 1, 8, 0, 6, 0, 2},
        {0, 8, 0, 0, 0, 0, 3, 0, 1},
        {0, 0, 0, 9, 0, 1, 0, 8, 0},
        {0, 0, 5, 0, 3, 0, 0, 0, 0},
        {0, 0, 0, 5, 0, 3, 0, 0, 0},
        {0, 9, 0, 4, 1, 0, 0, 0, 0},
        {0, 0, 3, 0, 0, 9, 0, 2, 0}
    };

    // Get the puzzle with corrected values
    int (*output_1)[9] = get_puzzle(input_1);

    // Print the original puzzle
    printf("Test 1 Input:\n");
    print_board(input_1);

    // Print the corrected puzzle
    printf("\nTest 1 ouput:\n");
    print_board(output_1);

    // TEST 2, check that a no solution puzzle, returns an empty puzzle
    int input_2[9][9] = {
        {6, 6, 0, 7, 0, 5, 1, 0, 4},
        {4, 5, 0, 0, 0, 0, 0, 9, 0},
        {3, 0, 0, 1, 8, 0, 6, 0, 2},
        {0, 8, 0, 0, 0, 0, 3, 0, 1},
        {0, 0, 0, 9, 0, 1, 0, 8, 0},
        {0, 0, 5, 0, 3, 0, 0, 0, 0},
        {0, 0, 0, 5, 0, 3, 0, 0, 0},
        {0, 9, 0, 4, 1, 0, 0, 0, 0},
        {0, 0, 3, 0, 0, 9, 0, 2, 0}
    };

    // Get the puzzle with corrected values
    int (*output_2)[9] = get_puzzle(input_2);

    // Print the original puzzle
    printf("Test 2 Input:\n");
    print_board(input_2);

    // Print the corrected puzzle
    printf("\nTest 2 ouput:\n");
    print_board(output_2);

    // TEST 3, check that a single solution puzzle, returns a the same puzzle that was input
    int input_3[9][9] = {
      {0, 6, 8, 7, 0, 5, 1, 0, 4},
      {4, 5, 0, 3, 0, 0, 0, 9, 8},
      {3, 0, 9, 1, 8, 0, 6, 5, 2},
      {0, 8, 0, 0, 0, 0, 3, 0, 1},
      {6, 0, 0, 9, 0, 1, 2, 8, 5},
      {0, 0, 5, 2, 3, 8, 0, 0, 0},
      {7, 0, 0, 5, 0, 3, 0, 0, 0},
      {0, 9, 6, 4, 1, 2, 5, 0, 3},
      {5, 1, 3, 8, 7, 9, 0, 2, 0}
  };

    // Get the puzzle with corrected values
    int (*output_3)[9] = get_puzzle(input_3);

    // Print the original puzzle
    printf("Test 3 Input:\n");
    print_board(input_3);

    // Print the corrected puzzle
    printf("\nTest 3 ouput:\n");
    print_board(output_3);

    return 0;
}
