## RSpace-CLI

RSpace-CLI is a command-line application to interact with RSpace ELN.

It is designed a supplement to the web interface for tasks such as:

* bulk upload or download of files
* bulk import of MSWord documents into native RSpace documents
* Querying audit-trail for activity
* Getting reports in JSON, tabular or CSV format
* Integrating with your data-management workflows

It is written in the Go programming language.

Run

    rspace eln --help

for commands and their arguments.

### Getting started

Once you have obtained or compiled the binary, you need to supply a configuration file.

Create a file called '.rspace' in your home folder and 2 lines with the URL of your RSpace and
your API key

    RSPACE_API_KEY=get_this_from_your_profile_page
    RSPACE_URL=https://myrspace.com/api/v1

[ ![Download](https://api.bintray.com/packages/ra22597/rspace-cli/rspace-cli/images/download.svg) ](https://bintray.com/ra22597/rspace-cli/rspace-cli/_latestVersion)

If you prefer you can add this information in any file and supply its filepath with the --config flag to each command, e.g.

    rspace eln listTree --config /path/to/myConfig.txt

