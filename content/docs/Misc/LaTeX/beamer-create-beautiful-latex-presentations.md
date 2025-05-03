---
weight: 999
url: "/Beamer_\\:_create_beautiful_LaTeX_presentations/"
title: "Beamer: Create Beautiful LaTeX Presentations"
description: "A comprehensive guide to using Beamer to create professional LaTeX presentations, with instructions for setup, formatting, and custom themes."
categories: ["LaTeX", "Presentation", "Development"]
date: "2013-09-04T15:58:00+02:00"
lastmod: "2013-09-04T15:58:00+02:00"
tags: ["latex", "beamer", "presentation", "slides", "documentation"]
toc: true
---

{{< table "table-hover table-striped" >}}
|||
|-|-|
| ![Beamer](/images/latex_logo.avif) ||
| **Software version** | 3.24-1 |
| **Operating System** | Debian 7 |
| **Website** | [Website](https://bitbucket.org/rivanvx/beamer/wiki/HomeBeamer) |
| **Last Update** | 04/09/2013 |
{{< /table >}}

## Introduction

Beamer[^1] is a LaTeX extension to create beautiful slides. The main goal of this tool is to stop spending hours and hours on the look but concentrate on the content of your slides.

Of course the first time you'll use it, you certainly spend some hours to understand how it works and it could be more if you never did LaTeX before. But don't worry it's possible as it nearly was my case.

There are some rules you need to understand if it's the first time you're using it:

1. If you need something quick...use LibreOffice or something else
2. If you have a lot of slides to do, Beamer is for you
3. If you absolutely need something clear, beautiful and that looks professional, use Beamer!

{{< alert context="info" text="I've made a custom theme with custom functionalities. Some of them require to load custom functions to work. I've added everything in that document as well" />}}

You can look at [a result example here](/pdf/beamer.pdf). If you want to grab the sources of that example, you can get it [from my Git](https://git.deimos.fr/?p=git_deimosfr.git;a=tree;f=others/beamer_template;hb=HEAD).

Now you're ready for the practice.

## Installation

To use LaTeX/Beamer, you don't need a lot of things, vim is enough! But I prefer having an IDE to help me on unknown syntax I'd like help, I'm using TexMaker. That's why we're going to install it as well:

```bash
aptitude install latex-beamer texmaker latexmk texlive-pictures
```

Latexmk is different from pdflatex (the default pdf generator of texmaker). In fact, when you generate your PDF on texmaker (F1 key shortcut), it will generate the summary and other things separately. But you may not have the definitive PDF at the first try, why?

Simply because this task needs to be performed several times to combine everything (content, summary....) to get the final document. Another solution exists to get a full document at once but will take more time to generate the final document. It's called latexmk and always should be used to generate final documents.

Create a folder containing images, scripts and other required documents:

```bash
mkdir -p beamer/{config,images}
touch beamer.tex
```

Then add the [theme](#custom-theme) and [custom functions](#custom-functions) in that folder as well.

We're now going to see how it works and how to write the slides.

## Start document

### Load libs

First of all you need to know that there is an order when you load libraries. That's why the first thing to load on your tex file is the libraries that you're going to use. Edit the beamer.tex file and add:

```latex
% slides format
\documentclass[aspectratio=43]{beamer}
%\documentclass[aspectratio=1610]{beamer}
%\documentclass[aspectratio=169]{beamer}

% encoding
\usepackage[english]{babel}
\usepackage[T1]{fontenc}
\usepackage[utf8]{inputenc}

% others
\usepackage{microtype}
\usepackage{lmodern}
\usepackage{hyperref}
\usepackage{bookmark}
\usepackage{textcomp}
\usepackage{graphicx}

% Symbols
\usepackage{amssymb}
\usepackage{latexsym}

% tikz / dia export
\usepackage{tikz}

% smart diagram
\usepackage{smartdiagram}

% theme
\useinnertheme{rounded}
%\usetheme{default}
\usetheme{deimos}
\usecolortheme{deimos}

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
```

- slides format: select your wished slides format (16/9, 4/3, 16/10)
- encoding: set encoding format for the document and the language
- others: libs to be able to insert hyperlinks, images...
- symbols: permit to insert LaTeX symbols like arrows
- tikz: will allow dia LaTeX exports to be included in the document
- smart diagram: make beautiful diagram easily with that lib (more info [here](https://mirrors.linsrv.net/tex-archive/graphics/pgf/contrib/smartdiagram/smartdiagram.pdf)[^2])
- theme: select theme and attributes (choose one [here](https://www.hartwork.org/beamer-theme-matrix/)[^3] or can create your custom theme [here](https://titilog.free.fr/)[^4])

### Document informations

The first things you need to add is variables containing document information:

```latex
\title{{My Title}}
\author{Your Name}
\institute{\href{http://www.mysite.com}{{\textcolor{blue}{www.mysite.com}}}}
\date{\today}
\logo{\includegraphics[height=5mm]{images/my_logo.png}\vspace{-7pt}}

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
```

Adapt all that information with yours. It will be reused in the future.

### Begin document and first frame

Now we're able to begin the document and set the welcome page:

```latex
\begin{document}

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

\begin{frame}[plain]

\titlepage
\pdfbookmark[2]{\inserttitle}{titre}
\begin{center}
\includegraphics[height=10mm]{images/my_logo.png}
\end{center}

\end{frame}

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
```

- \begin{document}: This is how to start the document
- \begin{frame}: This is how your start a frame (slide)
  - [plain]: specify that you don't want header and footer like on other slides
- {\inserttitle}{titre}: is how you insert information already defined above
- \includegraphics: permit to insert an image in a slide
  - [height=10mm]: set the image size

### Summary

A good thing would be to add a summary! To make it work automatically, add this:

```latex
\section*{Summary}

\pdfbookmark[2]{Summary}{Summary}
\begin{frame}[shrink]{Summary}
\tableofcontents[hideallsubsections]

\end{frame}

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
```

- \section: add a section that will appear in the summary (do not show it with the star)
- [shrink]: shrink automatically the frame in 2 frames if it's too big
- \tableofcontents: show the table of content
  - [hideallsubsections]: hide subsections

## Frames

### Sections and subsections

We're now going to create a section (corresponding to title 1):

```latex
\section{Title 1 (section)}

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
```

You can create subsection (title 2) like that:

```latex
\subsection{Title 2 (subsection)}
```

### Frames usages

You've seen how to create a frame! It's really easy. I suggest for all of your frames that you have this kind of blocks:

```latex
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

\begin{frame}{Title frame 2}

Content frame 1

\end{frame}

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

\begin{frame}{Title frame 2}

Content frame 2

\end{frame}

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
```

### Insert images

We already seen it but didn't see the alignment yet. You can decide to have a center alignment and including an image like this:

```latex
\begin{center}
\includegraphics[height=52mm]{images/image.jpg}
\end{center}
```

#### Vertical alignment

You can set vertical alignment of an image like this:

```latex
\begin{center}
\vspace{-25pt}
\includegraphics[height=52mm]{images/image.jpg}
\end{center}
```

### Insert Dia exported LaTeX images

In Dia, you can export to LaTeX and then import it like an image in your tex file. In Dia export in "Macros PGF LaTeX" format, then include it like that:

```latex
\begin{center}
\input{images/image.tex}
\end{center}
```

#### Resize imported Dia images

You can resize Dia images, in editing directly the .tex file and modifying this line at the beginning of the document:

```latex
\setlength{\du}{15\unitlength}
```

Change the number (here 15) by the desired size.

### Items list

To add an item list:

```latex
\begin{itemize}
\item item 1
\item item 2
\item item 3
\end{itemize}
```

### Symbols

You can insert symbols like arrows with '$' sign:

```latex
$\rightarrow$
$\leftarrow$
```

You can find a [list of symbols here](https://www.combinatorics.net/Resources/weblib/A.7/a7.html)[^5].

### Insert file content

I use it to add configuration file for example. Here is how to insert a file contained in the config folder:

```latex
\script{file.conf}{/path/to/file.conf}
```

This will keep format text without needing to insert escape chars...

### Insert command blocks

To insert a command block, you can do like that:

```latex
\begin{block}{mycommand}
\begin{lstlisting}
$ my_command argument1 argument2
\end{lstlisting}
\end{block}
```

This will keep format text without needing to insert escape chars...You will also need to insert the [fragile] element at the beginning of the frame:

```latex
\begin{frame}[fragile]{my_frame}
```

### Insert color blocks

#### Warning bloc

If you need to insert a warning block:

```latex
\begin{alertblock}{Warning}
Warning block
\end{alertblock}
```

#### Notes

If you want to insert a notes block:

```latex
\begin{exampleblock}{Notes}
Notes block
\end{exampleblock}
```

## Custom theme

Here is my theme with colors. There are normally enough comments to understand what each thing does:

```latex
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%
%             Deimos Beamer/LaTeX Presentation Color Theme
%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

\mode<presentation>

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% Document default font size
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

\usepackage{scrextend}
\changefontsizes{9.5pt}

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% Colors
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

\DefineNamedColor{named}{DeimosBlue}{RGB}{48,90,148}
\DefineNamedColor{named}{DeimosBlueLight}{RGB}{153,179,216}

\setbeamercolor{structure}{bg=black, fg=DeimosBlue}

\setbeamercolor{section in toc}{fg=black,bg=white}
\setbeamercolor{alerted text}{fg=DeimosBlue!80!DeimosBlueLight}
\setbeamercolor*{palette primary}{fg=DeimosBlue!60!black,bg=DeimosBlueLight!30!white}
\setbeamercolor*{palette secondary}{fg=DeimosBlue!70!black,bg=DeimosBlueLight!15!white}
\setbeamercolor*{palette tertiary}{bg=DeimosBlue!80!black,fg=DeimosBlueLight!10!white}
\setbeamercolor*{palette quaternary}{fg=DeimosBlue,bg=DeimosBlueLight!5!white}

\setbeamercolor*{sidebar}{fg=DeimosBlue,bg=DeimosBlueLight!15!white}

\setbeamercolor*{palette sidebar primary}{fg=DeimosBlue!10!black}
\setbeamercolor*{palette sidebar secondary}{fg=white}
\setbeamercolor*{palette sidebar tertiary}{fg=DeimosBlue!50!black}
\setbeamercolor*{palette sidebar quaternary}{fg=DeimosBlueLight!10!white}

%\setbeamercolor*{titlelike}{parent=palette primary}
\setbeamercolor{titlelike}{parent=palette primary,fg=DeimosBlue}
\setbeamercolor{frametitle}{bg=DeimosBlueLight!10!white}
\setbeamercolor{frametitle right}{bg=DeimosBlueLight!60!white}

\setbeamercolor*{separation line}{}
\setbeamercolor*{fine separation line}{}

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% Blocks
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

% Classic block
\setbeamercolor{block title}{fg=white,bg=DeimosBlue}
\setbeamercolor{block body}{fg=black,bg=DeimosBlueLight!30}
\setbeamerfont{block title}{size=\footnotesize}
\setbeamerfont{block body}{size=\footnotesize}
% Usage:
%\begin{block}{Normal block}
%  Texte du block \texttt{Normal block}
%\end{block}

% Warning block
\setbeamercolor{block title alerted}{fg=black,bg=red!70}
\setbeamercolor{block body alerted}{fg=black,bg=red!20}
% Usage:
%\begin{alertblock}{Warning block}
%  Texte du block \texttt{alertblock}
%\end{alertblock}

% Information block
\setbeamercolor{block title example}{fg=black,bg=yellow!70}
\setbeamercolor{block body example}{fg=black,bg=yellow!20}
% Usage:
%\begin{exampleblock}{Notes}
%  Exemple de block \texttt{Notes}
%\end{exampleblock}

\mode
<all>
```

## Custom functions

Here are my custom functions and document specificities. There are normally enough comments to understand what each thing does:

```latex
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
%
%                   Deimos Beamer/LaTeX Presentation Theme
%
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

\mode<presentation>

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% Presentation
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

% Footer
\useoutertheme[footline=authorinstitutetitle]{miniframes}
% Show frame number and total frame number
\expandafter\def\expandafter\insertshorttitle\expandafter{%
\insertshorttitle\hfill%
\insertframenumber\,/\,\inserttotalframenumber}

\setbeamercolor{separation line}{use=structure,bg=structure.fg!50!bg}

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% Items look
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

\setbeamertemplate{itemize items}[default]
\setbeamertemplate{enumerate items}[default]
\useinnertheme{rounded}

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% Redefine header
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

% Customize header / Enable bullet points if you remove it
% http://tex.stackexchange.com/questions/17288/is-it-possible-to-get-rid-of-the-bullets-in-the-miniframes-outer-theme

\setbeamertemplate{headline}{%
  % First line
  \begin{beamercolorbox}[ht=10pt,dp=1.125ex,%
      leftskip=.3cm,rightskip=.3cm plus1fil]{section in head/foot}
    \usebeamerfont{section in head/foot}\usebeamercolor[fg]{section in head/foot}%
    {\footnotesize \inserttitle : \insertsectionhead}
  \end{beamercolorbox}%
  % Separation
  \begin{beamercolorbox}[colsep=1.5pt]{middle separation line head}
  \end{beamercolorbox}
  % Second line
  \begin{beamercolorbox}[ht=0ex,dp=0ex,%
    leftskip=.3cm,rightskip=.3cm plus1fil]{subsection in head/foot}
    \usebeamerfont{subsection in head/foot}%\insertsubsectionhead
  \end{beamercolorbox}%
  % Separation
  \begin{beamercolorbox}[colsep=1.5pt]{lower separation line head}
  \end{beamercolorbox}
}

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% Margins
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% Margin size
\setbeamersize{text margin left=1em,text margin right=1em}
% Spaces between paragraphs
\setlength{\parskip}{5pt}

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% PDF should have good informations
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
\pdfinfo{
    /Title
    /Creator
    /Author
}

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% Subsection reminder
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
\AtBeginSubsection[]
{
  \frame<handout:0>
  {
    \frametitle{Plan}
    \tableofcontents[sectionstyle=show/hide,subsectionstyle=show/shaded/hide]
  }
}

%\defbeamertemplate*{Deimos}{miniframes theme}
%{%
%	\begin{beamercolorbox}[ht=2.25ex,dp=5ex]{section in head/foot}
%    	\insertnavigation{\paperwidth}
%		%\insertsectionnavigationhorizontal{\paperwidth}{}{\hfill\hfill}
%	\end{beamercolorbox}%
%}%

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% Table of content
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

% Toc with line numbers and dashes
\newcounter{sectionpage}
\long\def\beamer@section[#1]#2{%
  \beamer@savemode%
  \mode<all>%
  \ifbeamer@inlecture
    \refstepcounter{section}%
    \beamer@ifempty{#2}%
    {{\long\def\secname{#1}\long\def\lastsection{#1}}%
    {\global\advance\beamer@tocsectionnumber by 1\relax%
      \long\def\secname{#2}%
      \long\def\lastsection{#1}%
    \setcounter{sectionpage}{\insertframenumber}\stepcounter{sectionpage}%
      \addtocontents{toc}{\protect\beamer@sectionintoc{\the\c@section}{#2~\dotfill~\thesectionpage}{\the\c@page}{\the\c@part}%
        {\the\beamer@tocsectionnumber}}}%
    {\let\\=\relax\xdef\sectionlink{{Navigation\the\c@page}{\noexpand\secname}}}%
    \beamer@tempcount=\c@page\advance\beamer@tempcount by -1%
    \beamer@ifempty{#1}{}{}%
      \addtocontents{nav}{\protect\headcommand{\protect\sectionentry{\the\c@section}{#1}{\the\c@page}{\secname}{\the\c@part}}}%
      \addtocontents{nav}{\protect\headcommand{\protect\beamer@sectionpages{\the\beamer@sectionstartpage}{\the\beamer@tempcount}}}%
      \addtocontents{nav}{\protect\headcommand{\protect\beamer@subsectionpages{\the\beamer@subsectionstartpage}{\the\beamer@tempcount}}}%
    }%
    \beamer@sectionstartpage=\c@page%
    \beamer@subsectionstartpage=\c@page%
    \def\insertsection{\expandafter\hyperlink\sectionlink}%
    \def\insertsubsection{}%
    \def\insertsubsubsection{}%
    \def\insertsectionhead{\hyperlink{Navigation\the\c@page}{#1}}%
    \def\insertsubsectionhead{}%
    \def\insertsubsubsectionhead{}%
    \def\lastsubsection{}%
    \Hy@writebookmark{\the\c@section}{\secname}{Outline\the\c@part.\the\c@section}{2}{toc}%
    \hyper@anchorstart{Outline\the\c@part.\the\c@section}\hyper@anchorend%
    \beamer@ifempty{#2}{\beamer@atbeginsections}{\beamer@atbeginsection}%
  \fi%
  \beamer@resumemode}%

\def\beamer@subsection[#1]#2{%
  \beamer@savemode%
  \mode<all>%
  \ifbeamer@inlecture%
    \refstepcounter{subsection}%
    \beamer@ifempty{#2}{\long\def\subsecname{#1}\long\def\lastsubsection{#1}}
    {%
      \long\def\subsecname{#2}%
      \long\def\lastsubsection{#1}%
    \setcounter{sectionpage}{\insertframenumber}\stepcounter{sectionpage}%
      \addtocontents{toc}{\protect\beamer@subsectionintoc{\the\c@section}{\the\c@subsection}{#2~\dotfill~\thesectionpage}{\the\c@page}{\the\c@part}{\the\beamer@tocsectionnumber}}%
    }%
    \beamer@tempcount=\c@page\advance\beamer@tempcount by -1%
    \addtocontents{nav}{%
      \protect\headcommand{\protect\beamer@subsectionentry{\the\c@part}{\the\c@section}{\the\c@subsection}{\the\c@page}{\lastsubsection}}%
      \protect\headcommand{\protect\beamer@subsectionpages{\the\beamer@subsectionstartpage}{\the\beamer@tempcount}}%
    }%
    \beamer@subsectionstartpage=\c@page%
    \edef\subsectionlink{{Navigation\the\c@page}{\noexpand\subsecname}}%
    \def\insertsubsection{\expandafter\hyperlink\subsectionlink}%
    \def\insertsubsubsection{}%
    \def\insertsubsectionhead{\hyperlink{Navigation\the\c@page}{#1}}%
    \def\insertsubsubsectionhead{}%
    \Hy@writebookmark{\the\c@subsection}{#2}{Outline\the\c@part.\the\c@section.\the\c@subsection.\the\c@page}{3}{toc}%
    \hyper@anchorstart{Outline\the\c@part.\the\c@section.\the\c@subsection.\the\c@page}\hyper@anchorend%
    \beamer@ifempty{#2}{\beamer@atbeginsubsections}{\beamer@atbeginsubsection}%
  \fi%
  \beamer@resumemode}

% Triangle in summary
\setbeamertemplate{subsection in toc}{%
\leavevmode\leftskip=5.65ex%
  \llap{\raisebox{0.2ex}{\textcolor{structure}{$\blacktriangleright$}}\kern1ex}%
  \inserttocsubsection\par%
}

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% Other options
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

\DeclareOptionBeamer{compress}{\beamer@compresstrue}
\ProcessOptionsBeamer
% Remove navigation symbols
\setbeamertemplate{navigation symbols}{}

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% Listing
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

% Listings (used for formated text)
\RequirePackage{listingsutf8}

\lstset
{
    basicstyle      	= \ttfamily,
    inputencoding       = utf8/latin9,
    showstringspaces    = false,
    aboveskip			= -1pt,
    belowskip			= -1pt,
}

%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
% Functions
%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

\RequirePackage{fancyvrb}
\DefineVerbatimEnvironment{exemple}{Verbatim}{}
\DefineVerbatimEnvironment{intercom}{Verbatim}{commandchars=\\\{\}}

% Code syntax
% Usage : \script{fichier}
\newcommand{\script}[3][]{
        \smallskip
        \VerbatimInput
        [
                frame           = single ,
                framesep        = 0.5ex ,
                label           = {#3} ,
                obeytabs        = true ,
                fontsize        = \scriptsize,
                #1 ,
        ]
        {config/#2}
}

\mode
<all>
```

## References

[^1]: https://bitbucket.org/rivanvx/beamer/wiki/HomeBeamer
[^2]: http://mirrors.linsrv.net/tex-archive/graphics/pgf/contrib/smartdiagram/smartdiagram.pdf
[^3]: http://www.hartwork.org/beamer-theme-matrix/
[^4]: http://titilog.free.fr/
[^5]: http://www.combinatorics.net/Resources/weblib/A.7/a7.html
