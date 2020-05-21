# MeetingRoomScheduling
Coding Challenge for Delta Dental

## Setup

To enable initial setup and use of this application, first either clone or download this repository. Then in a terminal from
the same directory that `main.go` is in:, run:
```
go run main.go
```
The terminal you used to run this command should read information from the attached `input.txt` and `rooms.txt` files and
print the closest available applicable room to the console.

You're free to change the data in either file to test different room sizes, availability times, floors, and room numbers.
The only restriction is that the data in the .txt files should uniformly match the example data already present in them.

## How I solved the problem

I first copied the example code from the word document I recieved line-for-line, character-for-character, to ensure that
the code would be compatible with the testing input used by the author of that document. Afterward I converted the input
from the file into a string, then separated the data into a struct, which has the following qualities:

- Floor: (the floor the meeting room is on, this attribute is important because the goal of the assignment was to find the 
closest applicable meeting room to the floor that the team resides on)
- RoomNo: (this seemed less important, but I disconnected it from the floor number for better readability)
- Capacity: The amount the room can hold, it must be at least as large or greater than the team size
- Time Slots: each available time slot came in pairs, but I didn't bother with separate grouping, as the the pairs are 
ordered, so finding a time's significant other is as simple as increasing the array's index by one. (I did however convert
the time into base 100 for easier math operations [ex. 15 to 25, 30 to 50 etc...])
- Availability: I calculated the availability by subtracting an opening's end time by the opening's begining time. (This is 
where the math time conversions came in handy.

After extracting all of the data from the input file into a more usable struct, I used Binary Sort to search for the closet
floor that had a meeting room with adequate space and availability. Instead of using the middle of the Rooms array as the
insertion poin, I found the first available index that had the same Floor number as the request floor number. From there
it was up 1, down 2, up three, etc... . This way a comparison isn't needed between request floor numbers and available room
floor numbers.

## Optimal Testing

If I were to write tests for this function I would try to account for extreme test cases as well as failure conditions.
I would test for time slots too small or large, optimal meeting spots with floors both near and far away from the request 
floor, and I would also test for error handling.

## Improvements

Several improvements could be made to my method. For instance the binary sory assumes that the given `input.txt` file already
comes pre-sorted. (Although sorting the Rooms array by Floor number is trivial). Also my function doesn't account for syntax
errors or unexpected characters in the txt file (but it does have proper error handling for these conditions).


