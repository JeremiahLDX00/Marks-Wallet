# Documentation for Marks Wallet Package

## **_Things done during the checkpoint_**

<details><summary>Click here</summary>
<p>
  
#### Checkpoint submission promised
  <ul> 
    <li>Create database for Token and Transaction</li>
    <li>Create API GET functions for List all tokens and Search for Token</li>
    <li>Create API GET and POST functions for listing all transactions and making a transaction</li>
    <li>Post Github link for submission</li>
  </ul>
    
#### Checkpoint submission delivered
  <ul> 
    <li>Create database for Token and Transaction</li>
    <li>Create API GET functions for List all tokens and Search for Token</li>
    <li>Create API GET and POST functions for listing all transactions and making a transaction</li>
    <li>Post Github link for submission</li>
  </ul>
</p>
</details>

## **_Things done after checkpoint_**

<details><summary>Click here</summary>
<p>

#### After checkpoint tasks
  <ul> 
    <li>Clean up and reconfigure API for frontend</li>
    <li>Create frontend pages for marks wallet package</li>
    <li>Create docker containers for frontend, API and DB</li>
    <li>Create API documentation for body request and URLs</li>
    <li>Create Microservice Diagram</li>
    <li>Update readme for github</li>
    <li>Post Docker Links</li>
  </ul>
    
#### After checkpoint delivered
  <ul> 
    <li>Clean up and reconfigure API for frontend</li>
    <li>Create frontend pages for marks wallet package</li>
    <li>Create docker containers for frontend, API and DB</li>
    <li>Create API documentation for body request and URLs</li>
    <li>Create Microservice Diagram</li>
    <li>Update readme for github</li>
    <li>Post Docker Links</li>
  </ul>
</p>
</details>

## **Final Documentation**

<details><summary>Click here</summary>
<p>

#### Microservice Design
  
![image](https://user-images.githubusercontent.com/93190183/152667958-b7c3fc2e-f86e-4be2-be8b-45e222f2b1c7.png)
  
I have created this microservice using two APIs, one database and one frontend. The front end is categorized into 5 different HTML pages, the first page is the website home page where there's four different buttons for the user to interact with depending on which function they want to access:

<b>List all tokens</b>
<p>The list all tokens function gets the list of tokens that are available to the student that is in their wallet, after transactions have been done by the various microservices for the student, the function will pull all the tokens that the student has a balance with and display the tokens along with the amount of tokens remaining for the student.</p>

<b>Search for token</b>
<p>The search for token function references the Token type table that is in the database and retrieves the desired token according to the user's input, once they have entered an existing token, the system will then display the existing token to them.</p>

<b>List all transactions</b>
<p>The list all transactions function allows for the student to display all the transactions that has been done by them. Once they enter their student ID, the system will retrieve the data from the transactions table referencing their student ID and from there the student will then be able to see all the transactions they have made with other students respectively.</p>

<b>Make a transaction</b>
<p>The make a transaction function allows for the admin, tutor or student to create a transaction referencing their ID, the studentID they are sending to, the token ID they are using, the transaction type which is the description and the amount they are sending. once they have done submitted the transaction, if it succeeds there will be a message on the page that shows that their transaction has been created.</p>

<b>API Design</b>  
<p>The API's have been split according to their functionalities, which are shown in the above diagram as Token API and Transactions API respectively. This is so that the goal of being loosely coupled and for the microservice to be not too dependent on each other can be achieved. The Token API handles getting the list of tokens and searching for a token functions while the Transactions API handles listing the transactions as well as making a transaction.<p>
  
<b>Database Design</b>
<p>The Database has been designed in such a way that currently it has not been split for the various API's. The justification for this is such that other microservices will be able to easily access the same database to retrieve the data for the neccessary functions, all they need to do is to call the seperate API's respectively and they will be able to achieve getting the result from the various tables that have been made in the database.</p> 

<b>Docker</b>
<p>For this, docker has also been used to containerize the various functions. In this case, 4 containers have been created. 1 container for the frontend, 2 for the API's and 1 for the database. Once they have been containerized they can then be pushed into the server by cloning the repository and then composing up inside of the server to run the various microservices themselves.</p> 

![image](https://user-images.githubusercontent.com/93190183/152654863-dd92e036-59cf-44d4-8a12-7a3fc45ab692.png)

<p> This is the layout of the entire file. they have been categorized into the various folders and their various functions respectively, this is to achieve the goal of being loosely coupled to the best of my ability. This is so that they wouldn't be too dependent on each other later on and the services will then be able to operate on it's own.</p>
 
</p>
</details>

## **DockerHub Links**

<details><summary>Click here</summary>
<p>
  Frontend: https://hub.docker.com/r/jeremiahldx21/markswallet_frontend/tags<br>
  Database: https://hub.docker.com/r/jeremiahldx21/markswallet_markswalletdb/tags<br>
  Token API: https://hub.docker.com/r/jeremiahldx21/markswallet_tokens/tags<br>
  Transaction API: https://hub.docker.com/r/jeremiahldx21/markswallet_transactions/tags
</p>
</details>

## **Steps to Set up using GitHub**

<details><summary>Click here</summary>
<p>
  <ol>
   <li>Clone the repository using the URL that is in this website</li>
   <li>Run the docker-compose file and build using compose up</li>
   <li>check that the container and images have been created in your docker desktop</li>
   <li>Access the frontend container located in the container and run it via the browser</li>
   <li>You can now test the full package and see how it runs!</li>
  </ol>
</p>
</details>
