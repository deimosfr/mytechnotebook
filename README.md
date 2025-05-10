# My Tech Notebook (ex- Bloc Notes Info)

Here is what I have learned and discovered in the world of technology. I first use it for me as a Notebook, but I hope you will find it useful.

To run actions, you can use the [Just](https://github.com/casey/just) command.

# Local build

## Initiate the project

To initiate the project, you can use the following command:

```bash
just init
```

## Build and run the project

You just need to run:

```bash
just run
```

It's now accessible at [http://localhost:8000](http://localhost:8000).

## Code blocks highlights

Should be written this way:

```
```go hl_lines="2 3"
```go hl_lines="2-4 6"
```go hl_lines="2-4 6" linenums="1"
```

## Callouts

Example:

```
!!! note

    Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla et euismod
    nulla. Curabitur feugiat, tortor non consequat finibus, justo purus auctor
    massa, nec semper lorem quam in massa.
```

Variants:

```
!!! note "Phasellus posuere in sem ut cursus" # to name the callout
!!! note "" # without title
??? note # collapsed
???+ note # collapsable and expanded
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

Icons used can be:

* Material Design
* FontAwesome
* Octicons
* Simple Icons

Icon search can be done [here](https://squidfunk.github.io/mkdocs-material/reference/icons-emojis/)

# More info

* [Macros](https://mkdocs-macros-plugin.readthedocs.io/en/latest/)
