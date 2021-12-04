## Modou Niang

### Recipie Book
I have created a recipie book that will allow users to store recipies. 
Users can create new recipies by using the inputs at the top of the webpage. 
I have completed all required parts of the assignment, and everything is working 
as expected. One notable feature that is missing from my work is error handling. 
That being said, I am making the assumption that we will not have a rogue user, and 
all inputs will be passed in the correct data. New Recipies can be created with the 
input box at the top of the page. Recipies can be deleted with the delete button near 
each recipe. Recipies can also be edited with the edit button. Once the edit button 
is clicked, the user is free to edit the individual fields. This is simply done by 
clicking on the values and changing them. Once the user is satisfied with the changes, 
they can submit these these changes using the Submit button.

Assuming 3 backends, we can run this program using:
go run backend.go --listen 8090 --backend :8091,:8092, 
go run backend.go --listen 8091 --backend :8090,:8092, 
go run backend.go --listen 8092 --backend :8090,:8091. 

Our frontend can be ran with: go run frontend.go --listen 8080 --backend :8090,:8091,:8092. 

It is important to make sure the backends are started first so they can constantly be listening for any 
incoming requests. Once the backends are started, we sync them up by having them make synchronization 
requests with one another. From here we run Paxos everytime we make a change to the database - Create, Edit, 
and Delete. As of now, my application is working as expected and I have not discovered any issues.

Here is how my program handles the 5 test cases:
1. The frontend makes a call to all of the backends to check whether they are alive and requests for hte data. If theuy
are alive the request is fufilled, and our client continues to function.
2. Each backend contacts other backends as soon as it starts and compares the log lengths. If a backend has a 
longer log, we go with that backend as the source of truth.
3. After restarting all previously terminated replicas, disjoint sets of back end replicas will be forcibly
terminated and restarted, until all replicas have been terminated and restarted. Your application
should continue working normally during this process.
4. If this is case, we link the client to an error made and tell them to retry their request
5. This case is handled by syncing the logs everytime the app is started

I used the Paxos replication strategy for my application. I would have chosen raft, but I thought the leader selection process would have been much more difficult to implement and Paxos seemed pretty straightforward once it was understood. Each one of my backends has a log of messages that is appeneded to once the Paxos algorithm has been ran. The backends then use this Log to synchronize with one another. In order to prevent blocking while receiving messages, I run the receiving messages code on a seperate thread from the main thread. I also use a message to execute messages on a sepeerate thread, and the logs are synced with the other backends on a seperate thread.
Throughout the assignment, I have mainly used this website to https://people.cs.rutgers.edu/~pxk/417/notes/paxos.html help me better understand the Paxos implementation. 
