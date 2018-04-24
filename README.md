# Todo

Todo is a small CLI application for keeping track of your tasks. As an extra
twist, it will also add your tasks to Jira for visibility and tracability.

### Jira configuration
This step needs to be performed by a Jira administrator to create an Application link
for Todo. This only needs to be done once, and then you are good. Start by generating an RSA keypair

    openssl genrsa -out jira_privatekey.pem 1024
    openssl req -newkey rsa:1024 -x509 -key jira_privatekey.pem -out jira_publickey.cer -days 365
    openssl pkcs8 -topk8 -nocrypt -in jira_privatekey.pem -out jira_privatekey.pcks8
    openssl x509 -pubkey -noout -in jira_publickey.cer  > jira_publickey.pem

Now log in to the Jira administration page. Click **Application links** in the menu to the left.
Enter **Todo** in the input field and hit **Create new link**

Fill out the box something like this:

<img src="https://i.imgur.com/ysn7oRs.png" alt=":todo1" title=":todo1">

It does not really matter what you fill in here except the **Consumer key** and **Shared secret**

Proceed, and the next form you fill out like this:

<img src="https://i.imgur.com/ZINh2Jb.png" alt=":todo2" title=":todo2">

In the Public Key box, you paste the contents of **jira_publickey.pem** that you created earlier.
Now you are done. The only thing remaining is to distribute jira_privatekey.pem to any user wanting to try the tool

### Installation

#### Todo tool itself
You can download the latest release from the releases tab

If you prefer building yourself, clone this repo and build the tool using

    make build
This requires a working Go environment

#### RSA Key
Todo uses an RSA Key to sign its payloads so that Jira can verify it's an
authorized request. In order for this to work you need to place
**jira_privatekey.pem** in your **$HOME/.ssh** folder. If you don't want to / can't place it in that folder for
some reason you can specify its location with the **--privkey** flag

### Usage

    todo

Will bring up a menu of all available commands.

    todo help <command>
Will give you extended info about that specific command


#### Init
The first time you use the tool, you need to setup the Jira connection. You do
that by entering **todo init**
Follow the instructions on the screen and you will receieve an oauth token you
can use for all future requests.

    todo init
    Initalizing todo...

    Enter Jira URL: https://jira.company.company.com

    Visit https://jira.company.com/plugins/servlet/oauth/authorize?oauth_token=XXXXXXXXXXXXXXX give the tool access

    Are you done? (y/n): y

    Great! Authentication succesful! Let's proceed...

    Choose which project you want Todo to use:

    1) Project 1
    2) Project 2
    3) Project 3

    Pick a project: 2

    Project 'Project 2' chosen.

    Choose which status you want to use when marking an issue as completed:

    1) Backlog
    2) Selected for Development
    3) In Progress
    4) Done
    5) In Review
    6) Rejected

    Pick a status: 4

    Status 'Done' chosen.

    Initialization done! Enjoy the tool

#### List

    todo list
    +---+------+--------------------+---------------------+---------------------+-------------------------------------+
    | # | DONE |        TASK        |       CREATED       |      COMPLETED      |               URL                   |
    +---+------+--------------------+---------------------+---------------------+-------------------------------------+

Possible flags include **"-f open"** for listing only open tasks, and
**-f completed** for listing only completed tasks

#### Add

    $ todo add Testing todo tool
    $ Task added: 'Testing todo tool'
    $ todo list
    +---+------+--------------------+---------------------+-----------+-----------------------------------------------+
    | # | DONE |        TASK        |       CREATED       | COMPLETED |                      URL                      |
    +---+------+--------------------+---------------------+-----------+-----------------------------------------------+
    | 1 |      | testing Todo tool  | 2018-04-22 20:10:13 |           | https://jira.company.com/browse/PRJ-001       |
    +---+------+--------------------+---------------------+-----------+-----------------------------------------------+

#### Complete

    $ todo complete 1
    $ Completing task: 'Testing todo tool'
    $ todo list
    +---+------+--------------------+---------------------+---------------------+-----------------------------------------+
    | # | DONE |        TASK        |       CREATED       |      COMPLETED      |               URL                       |
    +---+------+--------------------+---------------------+---------------------+-----------------------------------------+
    | 1 |  X   | testing Todo tool  | 2018-04-22 20:10:13 | 2018-04-22 20:30:07 | https://jira.company.com/browse/PRJ-001 |
    +---+------+--------------------+---------------------+---------------------+-----------------------------------------+

#### Oops

    $ todo oops 1
    $ Re-openeing task: 'Testing todo tool'
    $ todo list
    +---+------+--------------------+---------------------+-----------+-----------------------------------------------+
    | # | DONE |        TASK        |       CREATED       | COMPLETED |                      URL                      |
    +---+------+--------------------+---------------------+-----------+-----------------------------------------------+
    | 1 |      | testing Todo tool  | 2018-04-22 20:10:13 |           | https://jira.company.com/browse/PRJ-001       |
    +---+------+--------------------+---------------------+-----------+-----------------------------------------------+

#### Del

    $ todo del 1
    $ Deleting task: 'Testing todo tool'
    $ todo list
    +---+------+--------------------+---------------------+---------------------+-------------------------------------+
    | # | DONE |        TASK        |       CREATED       |      COMPLETED      |               URL                   |
    +---+------+--------------------+---------------------+---------------------+-------------------------------------+

## Security concerns

Your Jira credentials are never used / logged / stored by the tool in any way. It
will however generate an oauth token that can be used to authenticate against Jira
as long as the token is valid. This token is encrypted with the RSA key mentioned above
before it gets stored on your disk.

So if an attacker gets access to this tools source code, your encrypted oauth
token, and the jira private RSA key, he could access Jira and perform actions in
your name. Keep this in mind. The author takes no responsibility for incorrect or reckless usage
of this tool.
