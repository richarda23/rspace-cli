## RSpace-CLI

RSpace-CLI is a command-line application to interact with RSpace ELN in a more convenient
and compact way than using the API directly.

It is designed a supplement to the web interface for tasks such as:

* bulk upload or download of files and folders
* bulk import of MSWord documents into native RSpace documents
* Querying audit-trail for activity
* Getting reports in JSON, tabular or CSV format
* Integrating cleanly into your data-management workflows
* Admin functions such as ad-hoc account creation

It is written in the Go programming language.

### Downloading

Signed executables for Linux, MacOSX and Windows (amd64 architecture) are available from Bintray:

[ ![Download](https://api.bintray.com/packages/ra22597/rspace-cli/rspace-cli/images/download.svg) ](https://bintray.com/ra22597/rspace-cli/rspace-cli/_latestVersion)

Download the latest version for your platform, rename to 'rspace' and check it works:

    rspace eln --help

to show commands and their arguments.

### Configuring

Next, you must supply a configuration file with your RSpace API credentials:

Create a file called '.rspace' in your home folder and add two lines with the URL of your RSpace and
your API key, like this:

    RSPACE_API_KEY=get_this_from_your_RSpace_profile_page
    RSPACE_URL=https://myrspace.com/api/v1

If you prefer, instead of the default '.rspace' file,  you can add this information to any file and supply its filepath with the --config flag to each command, e.g.

    rspace eln listTree --config /path/to/myConfig.txt

Using --config option is useful if you have more than one account (e.g. an admin account and a personal account)
