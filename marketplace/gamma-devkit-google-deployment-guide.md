# Google Cloud Launcher
This guide will ... deploying... testing and finally developing with...

## Deploying
First, we're going to deploy an Orbs development kit. 
The kit consists of:
* Gamma server - Our development local blockchain which stores our data and answers our calls)
* Gamma CLI - The CLI allows us to easily deploy contracts onto the blockchain and perform calls against our deployed contract.
* Prism - Orbs' block explorer, which provides visibility into blocks and transactions which have been registered on our network.

Deploying is easy with Cloud Launcher: Go to https://console.cloud.google.com/launcher/details/orbsltd-public/orbs-gamma-devkit

![](./images/step01.png)

Click "Launch on Compute Engine"

![](./images/step02.png)

You can take the default settings or customize them.  When complete click "Deploy"

![](./images/step03.png)

That's it!  Your development kit is now deploying.

## Inspecting the kit once it's running

When the deployment is complete, you should see:

![](./images/step04.png)

Note: Gamma, the development server, is configured to close a block once every 10 minutes or so, OR when at least 1 transaction is pushed through with an API call. 

At this point, you can open a web browser and view Prism to see block activity. Click on the "Visit the site" button.

![](./images/step05.png)

Once the browser has finished loading, it should look like:

![](./images/step06.png)

You should see a few blocks. That's great as it shows our development kit is working and closing blocks. 
<!--- I would write - four nodes are working in consensus and closing blocks.... -->
That means we can start to play around with it.

The next kind of interaction we can have with our development kit is to deploy a smart contract into the blockchain and interact with it.

## Developing with the deployed development kit

At the moment Prism and our Blockchain server are running and looks to be closing blocks,  which means it's healthy.

Next, we connect to the machine running our Gamma server so we can use Gamma CLI and deploy our first contract. We'll do so by clicking on "Connect":

![](./images/step07.png)

Google Cloud negotiates the SSH keys between our browser and the target machine in about 10 seconds. There are other means of accessing the machine via SSH. We are merely illustrating the most straightforward approach here for the sake of keeping things simple.

![](./images/step08.png)

Once inside the machine, our Terminal (which is in a new pop up window) should look like the following:

![](./images/step09.png)

Let's get some example contracts from [Orbs Contract SDK](https://github.com/orbs-network/orbs-contract-sdk) repository by running `git clone https://github.com/orbs-network/orbs-contract-sdk.git` within the terminal window.

(Tip: you might be prompted to approve the authenticity of github.com. Type `yes` into the Terminal if asked.)

![](./images/step10.png)

Let's navigate to the examples folder and into the Counter contract by typing `cd orbs-contract-sdk/go/examples/counter/` into the terminal window. The Counter contract is a very simple contract that initializes an empty state variable, which we can be increased with a call to the `add()` function. In addition, we can also access the current counter value by invoking a call to the contract's `get()` function. 

In the Terminal you should be located at the folder containing the file `contract.go`. Let's verify this by typing `ls -lh` into the Terminal to ensure we can see the file is there.

![](./images/step11.png)

Next, we can deploy the contract into our Gamma server by typing the following command: `gamma-cli deploy contract.go -name MyCounter` and wait for it to finish. Once finished, Gamma CLI will return a JSON output object representing the result of our transaction against the server. Which (if all is successful) should look similar to the following

![](./images/step12.png)

Please note the 2 most important properties in the JSON object:

* The `ExecutionResult` provides the most basic kind of proof that the request was successful or not. In this case, we can see that the call was successful.
* The `TransactionStatus` is very important. Since this is a blockchain, our transaction might not have been committed into any block yet which means it is not final, and therefore might be pending (or worse - not accepted at all into the ledger). In our case, we can see it's in the `COMMITTED` state so we can rest assure our contract is deployed successfully!

Congrats! You have successfully deployed your first smart contract! How easy was that?!?

Ok so what's next? If we look closely at our smart contract code, we can see it's initializing the counter value with a value of `0`.

![](./images/step13.png)

Since our contract is now deployed, let's try and "read" the current value of the counter to make sure we get a value of `0` back.

Gamma CLI uses a JSON file to represent the action details we're requesting against the blockchain. Our Counter example has JSON files ready for us to experiment with. The `get()` method doesn't need any arguments. The JSON descriptor file is quite easy to understand. Let's have a look at it:

Type into the terminal window `cat json/get.json` (cat is a command in Linux to view the content of a file)

![](./images/step14.png)

Note that this file is specifically asking to invoke the `get()` method of the `MyCounter` contract. That name is important as we used it previously to name our contract upon deployment. If you chose a different name, you have to edit this file before sending a request with Gamma CLI using this file.

Alright, so now that we have some sense about what gamma is using to send our request to the local blockchain server (Gamma), let's issue a "Query" call and invoke the `get()` command of our contract.

We do that by typing in the terminal `gamma-cli run-query json/get.json`

![](./images/step15.png)

Again, the returned output contains a bit of information. We've highlighted what's most important for us at this stage, and that's the current value of the counter state variable. We know it's supposed to be 0, but we wanted to verify that. Indeed as it appears by the highlighted red rectangle, the value is a `uint64` (a native number) of `0`.

All looks good so far!

Now let's make things interesting by running a transaction that adds '25` to that lonely `0` to increase the counter!

As before, let's view how the JSON file to describe this action is represented for Gamma CLI to use by typing `cat json/add-25.json`

![](./images/step16.png)

Notice how we have an argument in this call in the form of a `uint64` (native number) with a value of 25, which will be added to the counter upon running this transaction with Gamma CLI.

To make this happen, we issue the request with Gamma CLI by typing the following into the Terminal: `gamma-cli send-tx json/add-25.json`

and the input from this call should look something similar to the following:

![](./images/step17.png)

As you can see, the response seems to be successful! We've got an `ExecutionResult` of `SUCCESS`, which is excellent. It means our transaction went through <!--- add consensus talk? -->and was also `COMMITTED` into a block which is accepted by the network there being committed to the ledger. That's great news! We've modified the state of our contract now and updated the counter variable in its state to represent a different value (25).

Let's verify this is indeed the reality by issuing another `get()` call and assert that we get 25 and not 0 as before.

![](./images/step18.png)

Yes!
As expected we can see that the blockchain has returned us the current value of the counter, and it's now standing on 25 instead of the previous 0 that it once was. So everything is working!

## Next Steps

You've probably noticed that we've used `send-tx` and `run-query` on some of these calls. The `send-tx` is used for transactions which aim to change the state of a contract and perform and actual action. Whereas the `run-query` calls are usually "reading" information that is already present and doesn't require the heavy lifting of the consensus process to occur between the blockchain nodes as happens on a transaction.

To read more about this and learn more about blockchain in general and how to move forward into more complex examples using Orbs be sure to [visit our GitBook](https://orbs.gitbook.io/contract-sdk/).

To learn more about Orbs and it's mission feel free to [visit our website](https://www.orbs.com/).
