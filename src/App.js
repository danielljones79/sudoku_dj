// Start of Selection
import React, { useState, useEffect } from 'react';
// axios is a popular HTTP client for making API requests
import axios from 'axios';

// Use the backend service name in Docker, fallback to localhost for development
const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8081';

// Helper function to validate and transform board data from the new JSON format
const transformBoardData = (data) => {
  console.log('Raw response data:', data); // Debug log
  
  if (!data) {
    console.log('No data received');
    return [];
  }

  // Handle the new JSON structure with a "cells" object
  if (data.cells && typeof data.cells === 'object') {
    const cells = data.cells;
    const board = Array(9).fill(null).map(() => Array(9).fill(null));
    
    // Map cell IDs (e.g., "01", "02") to 2D array indices
    for (let key in cells) {
      const cell = cells[key];
      const index = parseInt(key, 10) - 1; // Convert "01" to 0-based index
      const row = Math.floor(index / 9);
      const col = index % 9;
      board[row][col] = {
        value: cell.value || 0,
        notes: Array.isArray(cell.notes) ? cell.notes : [],
        status: cell.status || ''
      };
    }
    console.log('Transformed cells object into 2D array');
    return board;
  }

  // Fallback for 2D array format
  if (Array.isArray(data) && data.length > 0 && Array.isArray(data[0])) {
    console.log('Data is already a 2D array');
    return data.map(row => 
      row.map(cell => {
        if (typeof cell === 'object') {
          return {
            value: cell.value || 0,
            notes: Array.isArray(cell.notes) ? cell.notes : [],
            status: cell.status || ''
          };
        } else {
          return {
            value: cell || 0,
            notes: [],
            status: ''
          };
        }
      })
    );
  }

  console.log('Could not transform data into a valid board');
  return [];
};

function App() {
  const [board, setBoard] = useState([]);
  const [difficulty, setDifficulty] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [puzzleId, setPuzzleId] = useState(null);
  const [availablePuzzles, setAvailablePuzzles] = useState([]);
  const [showPuzzleList, setShowPuzzleList] = useState(false);
  const [notesMode, setNotesMode] = useState(false);
  const [selectedCell, setSelectedCell] = useState(null);
  const [highlightedValue, setHighlightedValue] = useState(null);
  const [message, setMessage] = useState(null); // New state for temporary messages
  const [deleteConfirmation, setDeleteConfirmation] = useState(null); // State for delete confirmation

  // Fetch the list of available puzzles when component mounts
  // and load the last puzzle or generate a new one
  useEffect(() => {
    // Use a flag to ensure we only run this effect once
    let isMounted = true;
    
    const initializeApp = async () => {
      try {
        setLoading(true);
        setError(null);
        console.log('Fetching available puzzles on initial load');
        const response = await axios.get(`${API_BASE_URL}/sudoku`, {
          headers: {
            'Accept': 'application/json'
          }
        });
        
        if (!isMounted) return;
        
        console.log('Initial puzzles list response:', response);
        if (response.data && Array.isArray(response.data)) {
          // Log difficulty values for debugging
          response.data.forEach(puzzle => {
            console.log(`Puzzle ${puzzle.shortId || puzzle.uuid.substring(0, 8)} - Difficulty: ${puzzle.difficulty}`);
          });
          
          setAvailablePuzzles(response.data);
          
          if (response.data.length > 0) {
            // Load the most recent puzzle (the backend already sorts by newest first)
            const mostRecentPuzzle = response.data[0];
            console.log('Loading most recent puzzle:', mostRecentPuzzle.uuid);
            await loadPuzzle(mostRecentPuzzle.uuid);
            showMessage('Puzzles loaded successfully', 'success'); // Example usage
          } else {
            // If no puzzles are available, generate a new one
            console.log('No puzzles available, generating new puzzle');
            await generatePuzzle();
          }
        } else {
          // If response is invalid, generate a new puzzle
          console.log('Invalid puzzle list response, generating new puzzle');
          await generatePuzzle();
        }
      } catch (err) {
        if (!isMounted) return;
        setError('Failed to load puzzle list');
        console.error('Initial fetch error:', err);
        
        // If we can't load the puzzle list, generate a new puzzle
        console.log('Error fetching puzzles, generating new puzzle');
        await generatePuzzle();
        showMessage('Failed to load puzzles', 'error'); // Example error message
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    };
    
    initializeApp();
    
    // Cleanup function to prevent state updates after unmount
    return () => {
      isMounted = false;
    };
  // Using an empty dependency array but these functions are stable references
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Add event listener to clear selected cell when clicking outside the board
  useEffect(() => {
    const handleOutsideClick = (e) => {
      // Check if the click is outside the board
      if (!e.target.closest('.board')) {
        setSelectedCell(null);
        setHighlightedValue(null);  // Clear highlighted value when clicking outside
      }
    };

    document.addEventListener('click', handleOutsideClick);
    return () => {
      document.removeEventListener('click', handleOutsideClick);
    };
  }, []);

  // Add keyboard event listener for notes mode toggle and number input
  useEffect(() => {
    const handleKeyDown = (e) => {
      // Toggle notes mode when 'n' key is pressed
      if (e.key === 'n' && !e.ctrlKey && !e.altKey && !e.metaKey) {
        setNotesMode(prevMode => !prevMode);
        return;
      }

      // Handle number key presses when a cell is selected
      if (selectedCell && /^[1-9]$/.test(e.key)) {
        const { row, col } = selectedCell;
        const numValue = parseInt(e.key, 10);
        
        // Don't allow changes to static cells
        if (board[row][col].status === 's') return;
        
        if (notesMode) {
          // Toggle note in notes mode
          handleNoteToggle(row, col, numValue);
        } else {
          // Set cell value in normal mode (even if it has notes)
          handleCellValue(row, col, numValue);
          
          // Update highlighted value
          setHighlightedValue(numValue);
        }
      }

      // Clear cell and highlighted value with delete or backspace
      if (selectedCell && (e.key === 'Delete' || e.key === 'Backspace')) {
        const { row, col } = selectedCell;
        
        // Don't allow changes to static cells
        if (board[row][col].status === 's') return;
        
        // Clear the cell value
        handleCellValue(row, col, 0);
        
        // Clear highlighted value
        setHighlightedValue(null);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => {
      window.removeEventListener('keydown', handleKeyDown);
    };
  }, [selectedCell, notesMode, board]);

  // Function to fetch all available puzzles
  const fetchPuzzleList = async () => {
    try {
      setLoading(true);
      setError(null);
      console.log('Fetching available puzzles');
      const response = await axios.get(`${API_BASE_URL}/sudoku`, {
        headers: {
          'Accept': 'application/json'
        }
      });
      console.log('Puzzles list response:', response);
      if (response.data && Array.isArray(response.data)) {
        // Log difficulty values for debugging
        response.data.forEach(puzzle => {
          console.log(`Puzzle ${puzzle.shortId || puzzle.uuid.substring(0, 8)} - Difficulty: ${puzzle.difficulty}`);
        });
        setAvailablePuzzles(response.data);
      }
    } catch (err) {
      setError('Failed to load puzzle list');
      console.error('Fetch error:', err);
      console.error('Error response:', err.response);
    } finally {
      setLoading(false);
    }
  };

  const generatePuzzle = async () => {
    try {
      setLoading(true);
      setError(null);
      console.log('Generating puzzle with difficulty:', difficulty);
      const response = await axios.post(`${API_BASE_URL}/sudoku`, null, {
        params: {
          difficulty: difficulty
        },
        headers: {
          'Accept': 'application/json'
        }
      });
      console.log('Generate response:', response);
      const transformedBoard = transformBoardData(response.data);
      setBoard(transformedBoard);
      
      // Store the puzzle UUID for later use
      if (response.data.uuid) {
        setPuzzleId(response.data.uuid);
        
        // Add the new puzzle to availablePuzzles
        const newPuzzle = {
          uuid: response.data.uuid,
          shortId: response.data.uuid.substring(0, 8),
          difficulty: difficulty,
          date: new Date().toISOString()
        };
        setAvailablePuzzles(prev => [newPuzzle, ...prev]);
      }

      // Do not refresh the puzzle list here - can cause a loop
      // fetchPuzzleList();
    } catch (err) {
      setError('Failed to generate puzzle');
      console.error('Generate error:', err);
      console.error('Error response:', err.response);
    } finally {
      setLoading(false);
    }
  };

  // Function to load a specific puzzle by UUID
  const loadPuzzle = async (uuid) => {
    try {
      setLoading(true);
      setError(null);
      console.log('Loading puzzle with UUID:', uuid);
      const response = await axios.get(`${API_BASE_URL}/sudoku/${uuid}`, {
        headers: {
          'Accept': 'application/json'
        }
      });
      console.log('Load puzzle response:', response);
      const transformedBoard = transformBoardData(response.data);
      setBoard(transformedBoard);
      setPuzzleId(uuid);
      setShowPuzzleList(false);
    } catch (err) {
      setError('Failed to load puzzle');
      console.error('Load error:', err);
      console.error('Error response:', err.response);
    } finally {
      setLoading(false);
    }
  };

  // Function to save current puzzle state
  const savePuzzle = async () => {
    if (!board.length || !puzzleId) {
      setError('No puzzle to save');
      return;
    }
    
    try {
      setLoading(true);
      setError(null);
      console.log('Saving puzzle with UUID:', puzzleId);
      
      // Convert the 2D board array to the format expected by the backend
      const requestData = {
        uuid: puzzleId,
        cells: {},
        difficulty: difficulty // Include the current difficulty setting
      };
      
      // Convert board to cells format
      board.forEach((row, rowIdx) => {
        row.forEach((cell, colIdx) => {
          const posKey = String((rowIdx * 9 + colIdx + 1)).padStart(2, '0');
          requestData.cells[posKey] = {
            value: cell.value,
            notes: cell.notes || [],
            status: cell.status || (cell.value ? 'u' : '') // Mark user entries as 'u'
          };
        });
      });
      
      const response = await axios.put(`${API_BASE_URL}/sudoku/${puzzleId}`, requestData, {
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      });
      console.log('Save response:', response);
      
      // Update board with any changes from server
      const transformedBoard = transformBoardData(response.data);
      setBoard(transformedBoard);
      
      // Show success message
      showMessage('Puzzles saved successfully', 'success');
    } catch (err) {
      showMessage('Puzzles save failed', 'error');
      setError('Failed to save puzzle');
      console.error('Save error:', err);
      console.error('Error response:', err.response);
    } finally {
      setLoading(false);
    }
  };

  const validatePuzzle = async () => {
    if (!board.length) {
      setError('No puzzle to validate');
      return;
    }
    
    try {
      setLoading(true);
      setError(null);
      console.log('Validating puzzle:', board);
      
      // Convert the 2D board array back to the format expected by the backend
      const requestData = {
        uuid: puzzleId,
        cells: {},
        difficulty: difficulty // Include the current difficulty setting
      };
      
      // Convert board to cells format
      board.forEach((row, rowIdx) => {
        row.forEach((cell, colIdx) => {
          const posKey = String((rowIdx * 9 + colIdx + 1)).padStart(2, '0');
          requestData.cells[posKey] = {
            value: cell.value,
            notes: cell.notes || [],
            status: cell.status || (cell.value ? 'u' : '') // Mark user entries as 'u'
          };
        });
      });
      
      const response = await axios.post(`${API_BASE_URL}/sudoku/${puzzleId}`, requestData, {
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        }
      });
      console.log('Validate response:', response);
      const transformedBoard = transformBoardData(response.data);
      setBoard(transformedBoard);
    } catch (err) {
      setError('Failed to validate puzzle');
      console.error('Validate error:', err);
      console.error('Error response:', err.response);
    } finally {
      setLoading(false);
    }
  };

  // Function to handle toggling a note value
  const handleNoteToggle = (rowIdx, colIdx, noteValue) => {
    const newBoard = board.map((row, i) =>
      row.map((cell, j) => {
        if (i === rowIdx && j === colIdx) {
          const notes = [...(cell.notes || [])];
          const noteIndex = notes.indexOf(noteValue);
          
          if (noteIndex >= 0) {
            // Remove the note if it already exists
            notes.splice(noteIndex, 1);
          } else {
            // Add the note if it doesn't exist
            notes.push(noteValue);
            notes.sort((a, b) => a - b); // Keep notes sorted
          }
          
          return { 
            ...cell, 
            value: 0, // Clear the cell value when adding notes
            notes: notes,
            status: cell.status === 's' ? 's' : '' // Preserve system cells
          };
        }
        return cell;
      })
    );
    setBoard(newBoard);
  };

  // Function to set a cell value (clearing any notes)
  const handleCellValue = (rowIdx, colIdx, value) => {
    const newBoard = board.map((row, i) =>
      row.map((cell, j) => {
        if (i === rowIdx && j === colIdx) {
          return { 
            ...cell, 
            value: value, 
            notes: [], // Clear notes when setting a value
            status: cell.status === 's' ? 's' : (value ? 'u' : '') // Preserve system cells, mark user entries
          };
        }
        return cell;
      })
    );
    setBoard(newBoard);
  };

  // Function to handle cell selection
  const handleCellSelect = (rowIdx, colIdx) => {
    // Don't select static cells
    if (board[rowIdx][colIdx].status === 's') {
      // Still allow highlighting the value for static cells
      const cellValue = board[rowIdx][colIdx].value;
      if (cellValue > 0) {
        setHighlightedValue(cellValue);
      }
      // But don't make them the selected editable cell
      setSelectedCell(null);
      return;
    }
    
    setSelectedCell({ row: rowIdx, col: colIdx });
  };

  // Handle cell input changes from input field
  const handleCellChange = (rowIdx, colIdx, value) => {
    // Don't allow changes to static cells
    if (board[rowIdx][colIdx].status === 's') return;

    const newValue = value === '' ? 0 : parseInt(value, 10);
    if (isNaN(newValue) || newValue < 0 || newValue > 9) return;

    if (notesMode && newValue > 0) {
      // In notes mode, toggle the number in the notes array
      handleNoteToggle(rowIdx, colIdx, newValue);
    } else {
      // In normal mode, update the cell value and clear notes
      handleCellValue(rowIdx, colIdx, newValue);
      
      // Update highlighted value if we change the value of the selected cell
      if (newValue > 0) {
        setHighlightedValue(newValue);
      } else {
        setHighlightedValue(null);
      }
    }
  };

  // Handle adding/removing a specific note when clicked in notes mode
  const handleNoteClick = (rowIdx, colIdx, noteValue) => {
    // Don't allow changes to static cells
    if (board[rowIdx][colIdx].status === 's') return;

    // Only handle note clicks when in notes mode
    if (!notesMode) {
      // In normal mode, clicking on a note cell sets the value to that note
      handleCellValue(rowIdx, colIdx, noteValue);
      return;
    }

    handleNoteToggle(rowIdx, colIdx, noteValue);
  };

  // Function to handle cell click
  const handleCellClick = (rowIdx, colIdx) => {
    // Set this as the selected cell
    handleCellSelect(rowIdx, colIdx);
    
    // Highlight cells with the same value (if value > 0)
    const cellValue = board[rowIdx][colIdx].value;
    if (cellValue > 0) {
      setHighlightedValue(cellValue);
    } else {
      setHighlightedValue(null);
    }
  };

  // Function to show a message for 3 seconds
  const showMessage = (text, type = 'success') => {
    setMessage({ text, type });
    setTimeout(() => setMessage(null), 3000); // Clear after 3 seconds
  };

  // Function to delete a puzzle
  const deletePuzzle = async (uuid) => {
    try {
      setLoading(true);
      setError(null);
      console.log('Deleting puzzle with UUID:', uuid);
      const response = await axios.delete(`${API_BASE_URL}/sudoku/${uuid}`, {
        headers: {
          'Accept': 'application/json'
        }
      });
      console.log('Delete puzzle response:', response);
      
      // Remove deleted puzzle from availablePuzzles
      setAvailablePuzzles(prev => prev.filter(puzzle => puzzle.uuid !== uuid));
      
      // If the deleted puzzle is the current one, generate a new one
      if (uuid === puzzleId) {
        console.log('Deleted current puzzle, generating new one');
        setPuzzleId(null);
        await generatePuzzle();
      }
      
      // Show success message
      showMessage('Puzzle deleted successfully', 'success');
      
      // Clear delete confirmation
      setDeleteConfirmation(null);
    } catch (err) {
      setError('Failed to delete puzzle');
      console.error('Delete error:', err);
      console.error('Error response:', err.response);
      showMessage('Failed to delete puzzle', 'error');
    } finally {
      setLoading(false);
    }
  };

  // Function to handle delete confirmation
  const confirmDeletePuzzle = (uuid, shortId) => {
    setDeleteConfirmation({ uuid, shortId });
  };

  // Function to cancel delete confirmation
  const cancelDeletePuzzle = () => {
    setDeleteConfirmation(null);
  };

  return (
    <div className="App">
      <h1>Sudoku DJ</h1>
      <div className="controls">
        <label className="difficulty-control">
          Difficulty:
          <input
            type="number"
            min="1"
            max="9"
            value={difficulty}
            onChange={(e) => setDifficulty(Number(e.target.value))}
          />
        </label>
        <button onClick={generatePuzzle} disabled={loading} className="generate-btn">
          Generate Puzzle
        </button>
        <div className={`notes-toggle ${notesMode ? 'active' : ''}`}>
          <label className="toggle-switch">
            <input
              type="checkbox"
              checked={notesMode}
              onChange={() => setNotesMode(!notesMode)}
            />
            <span className="toggle-slider"></span>
          </label>
          <span className="toggle-label">
            Notes Mode <strong>{notesMode ? "(ON)" : "(OFF)"}</strong>
            <span className="keyboard-shortcut"> (Press 'n')</span>
          </span>
        </div>
      </div>
      {error && <div className="error">{error}</div>}
      {loading && <div className="loading">Loading...</div>}
      
      {/* Delete Confirmation Dialog */}
      {deleteConfirmation && (
        <div className="delete-confirmation">
          <p>Are you sure you want to delete puzzle <strong>{deleteConfirmation.shortId}</strong>?</p>
          <div className="delete-buttons">
            <button 
              className="delete-confirm" 
              onClick={() => deletePuzzle(deleteConfirmation.uuid)}
              disabled={loading}
            >
              Yes, Delete
            </button>
            <button 
              className="delete-cancel" 
              onClick={cancelDeletePuzzle}
              disabled={loading}
            >
              Cancel
            </button>
          </div>
        </div>
      )}
      
      {/* Main content area with two-column layout */}
      <div className="main-content">
        {/* Sudoku Board Section */}
        <div className="board-section">
          {board.length > 0 && (
            <div className="board-container">
              <div className="board">
                {board.map((row, i) => (
                  <div key={i} className="row">
                    {row.map((cell, j) => (
                      <div 
                        key={j} 
                        className={`cell ${cell.status === 's' ? 'static' : ''} 
                                   ${cell.status === 'c' ? 'correct' : ''} 
                                   ${cell.status === 'w' ? 'wrong' : ''} 
                                   ${cell.notes && cell.notes.length > 0 ? 'has-notes' : ''}
                                   ${selectedCell && selectedCell.row === i && selectedCell.col === j ? 'selected' : ''}
                                   ${highlightedValue && cell.value === highlightedValue ? 'highlighted' : ''}`}
                        onClick={() => handleCellClick(i, j)}
                      >
                        {cell.notes && cell.notes.length > 0 ? (
                          <div className="notes-grid">
                            {[1, 2, 3, 4, 5, 6, 7, 8, 9].map(num => (
                              <div 
                                key={num} 
                                className={`note-cell ${cell.notes.includes(num) ? 'note-visible' : ''}`}
                                onClick={(e) => {
                                  e.stopPropagation();
                                  handleNoteClick(i, j, num);
                                }}
                              >
                                {cell.notes.includes(num) ? num : ''}
                              </div>
                            ))}
                          </div>
                        ) : (
                          <input
                            type="text"
                            value={cell.value === 0 ? '' : cell.value}
                            onChange={(e) => handleCellChange(i, j, e.target.value)}
                            maxLength="1"
                            disabled={cell.status === 's'} // Disable input for static cells
                            onFocus={() => handleCellSelect(i, j)}
                          />
                        )}
                      </div>
                    ))}
                  </div>
                ))}
              </div>
            </div>
          )}
          
          {/* Validate and Save buttons below the puzzle grid */}
          <div className="board-controls">
            <button onClick={validatePuzzle} disabled={loading || !board.length} className="validate-btn">
              Validate Puzzle
            </button>
            <button onClick={savePuzzle} disabled={loading || !board.length} className="save-btn">
              Save Progress
            </button>
          </div>
        </div>
        
        {/* Puzzle List Section - always shown */}
        <div className="puzzle-list-section">
          <div className="puzzle-list">
            <h2>Available Puzzles</h2>
            {/* Call fetchPuzzleList on mount to ensure latest puzzles */}
            {useEffect(() => { fetchPuzzleList(); }, [])}
            {availablePuzzles.length === 0 ? (
              <p>No puzzles available. Generate one to get started!</p>
            ) : (
              <ul>
                {availablePuzzles.map((puzzle) => (
                  <li key={puzzle.uuid}>
                    <div className="puzzle-item">
                      <button 
                        className="load-puzzle-btn"
                        onClick={() => loadPuzzle(puzzle.uuid)}
                      >
                        <div className="puzzle-info">
                          <span className="puzzle-id">{puzzle.shortId || puzzle.uuid.substring(0, 8)}</span>
                          <span className="puzzle-date">{new Date(puzzle.date).toLocaleDateString()}</span>
                          <span className="puzzle-difficulty">
                            Difficulty: {puzzle.difficulty}
                            <div className="difficulty-bar">
                              <div 
                                className="difficulty-level" 
                                style={{width: `${puzzle.difficulty * 11}%`}}
                              ></div>
                            </div>
                          </span>
                        </div>
                      </button>
                      <button 
                        className="delete-puzzle-btn"
                        onClick={() => confirmDeletePuzzle(puzzle.uuid, puzzle.shortId || puzzle.uuid.substring(0, 8))}
                        title="Delete puzzle"
                      >
                        ×
                      </button>
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>
      </div>
      
      {message && (
        <div 
          className="temp-message"
          style={{
            color: message.type === 'error' ? 'red' : 'green',
            position: 'fixed',
            top: '20px',
            left: '50%',
            transform: 'translateX(-50%)',
            padding: '10px 20px',
            backgroundColor: '#fff',
            borderRadius: '4px',
            boxShadow: '0 2px 5px rgba(0,0,0,0.2)',
            zIndex: 1000
          }}
        >
          {message.text}
        </div>
      )}
    </div>
  );
}

export default App; 