# go-env
Personal environment variable allocator, etc, just testing go packages really. 

## How to use

Create an EnvironmentContext variable, and call the .Setup() method.

Will look for a ".env.*" file, written in json, example below:

.env.local
{
  "myVariable": "example"
}

Methods available for the Context are:

Get(string) -> Returns the value of an environment variable, using os.lookup.
GetVariables() -> Returns variables stored in the Environment, this probably isn't good practice, can get rid of this as it doesn't contribute much to anything.
GetEnvironment() -> Returns the current environment, (i.e .env.local would return "local").
GetKey() -> Method of the EnvironmentVariable type, returns name of the variable.

Very simple, testing go package creation.
