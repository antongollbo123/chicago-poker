# Chicago Poker

A multiplayer command-line Chicago Poker game playable via netcat!

## What is Chicago Poker?

Chicago Poker is a variant that combines poker and trick-taking:
- **Poker Rounds**: Players receive 5 cards, can toss and redraw, then the best hand wins points
- **Trick Rounds**: After 3 poker rounds, play 5 tricks (like bridge/hearts) for bonus points
- **Win Condition**: First player to 50 points wins!

## Quick Start

### 1. Build the game
```bash
go build ./cmd/chicago-poker
```

### 2. Start the server
```bash
./start.sh
```

### 3. Connect players (in separate terminals)
```bash
nc localhost 8080
```

Game starts automatically when 2 players connect!

## Scripts

- `./start.sh` - Start the server with live logs
- `./test.sh` - Run automated tests to verify the game works
- `./kill.sh` - Kill any running servers

## How to Play

**Poker Round:**
- View your hand (indexed 0-4)
- Enter indices of cards to toss (e.g., `0 2 4`) or press Enter to keep all
- Best hand wins points (1-8 based on poker rank)

**Trick Round:**
- Enter a single card index to play (e.g., `0`)
- Must follow suit if possible
- Highest card of lead suit wins the trick
- Winner of final trick gets 3 points

## Project Structure

```
cmd/chicago-poker/        Main entry point
internal/
  ├── gameNetwork/        Network game logic & server
  ├── game/              Hand evaluation & core rules
  ├── deck/              Deck management
  ├── player/            Player data structure
  └── utils/             Utility functions
pkg/cards/              Card data structures
```

## TODO

### Short term
- ✅ ~~Fix loop logic in trick round~~
- ✅ ~~Make game playable via netcat~~
- Add support for 3-4 players
- Better error messages

### Long term
- Editable config file
- Cloud deployment
- Web interface
- Tournament mode

