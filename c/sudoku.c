#include <stdio.h>
#include <stdbool.h>

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
                    printf("Found Naked Single %d,%d = %d\n", i, j, possible_values[0]);
                    board[i][j] = possible_values[0];
                }
            }
        }
    }
}
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

bool solve_sudoku(int *board) {
    int (*board_2d)[9] = (int (*)[9])board;
    fill_naked_singles((int (*)[9])board);
    fill_unique_candidates((int (*)[9])board);
    return solve_sudoku_recursive((int (*)[9])board);
}

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

int count_solutions(int *board) {
    int (*board_2d)[9] = (int (*)[9])board;
    return count_solutions_recursive(board_2d);
}
