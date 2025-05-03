# My Tech Notebook (ex- Bloc Notes Info)

Here is what I have learned and discovered in the world of technology. I first use it for me as a Notebook, but I hope you will find it useful.

To run actions, you can use the [task](https://taskfile.dev/) command.

# Local build

## Initiate the project

To initiate the project, you can use the following command:

```bash
task init
```

## Build and run the project

You just need to run:

```bash
task run
```

It's now accessible at [http://localhost:1313](http://localhost:1313).

# Write documentation

## Code blocks highlights

Should be written this way:

```
\`\`\`go {linenos=table,hl_lines=[3]}
\`\`\`go {linenos=table,hl_lines=[3],anchorlinenos=true}
\`\`\`go {linenos=table,hl_lines=[3,"5-7"],linenostart=199,anchorlinenos=true}
```

## Callouts

You can use the following callouts:

```

{{< alert context="info" text="" />}}

```

Possible contexts:

- success
- danger
- warning
- primary
- light
- dark

Or for a more complex case with html:

```

{{% alert icon="🛒" context="info" %}}
xxx
{{% /alert %}}

```

## Internal references

Example:

```
{{< ref "docs/Perl/_index.md" >}}
[xxx]({{< ref "docs/coding/php/" >}}).
[Perl]({{< ref "docs/Coding/Perl/introduction_to_perl.md#chomp" >}})
[SSLH method]({{< ref "docs/Servers/File sharing/SFTP and FTP/sslh_multiplexing_ssl_and_ssh_connections_on_the_same_port.md" >}})
```

## Tabs

```
{{< tabs tabTotal="3">}}
{{% tab tabName="Windows" %}}

**Windows Content**

Example content specific to **Windows** operating systems

{{% /tab %}}
{{% tab tabName="MacOS" %}}

**MacOS Content**

Example content specific to **Mac** operating systems

{{% /tab %}}
{{% tab tabName="Linux" %}}

**Linux Content**

Example content specific to **Linux** operating systems

{{% /tab %}}
{{< /tabs >}}
```

## icons

Icons can be found here:

- https://fonts.google.com/icons
- https://simpleicons.org/

# More info

https://lotusdocs.dev/
