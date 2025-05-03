---
weight: 999
url: "/La_programmation_orientée_objet_en_Perl/"
title: "Object-Oriented Programming in Perl"
description: "A guide to understanding and implementing Object-Oriented Programming in Perl, including basic concepts and practical examples"
categories: 
  - "Linux"
  - "Development"
date: "2012-10-13T20:15:00+02:00"
lastmod: "2012-10-13T20:15:00+02:00"
tags: 
  - "Perl"
  - "Development"
  - "Programming"
  - "Object-Oriented"
  - "Special pages"
  - "View source"
  - "Network"
  - "Printable version"
  - "Perl Website"
  - "cd ~"
  - "Servers"
  - "Windows"
  - "What links here"
toc: true
---

![Perl](/images/perl_logo.avif)

{{< table "table-hover table-striped" >}}
|||
|-|-|
| **Software version** | 5.10 |
| **Website** | [Perl Website](https://www.perl.org) |
| **Last Update** | 13/10/2012 |
{{< /table >}}

## Introduction

Object-Oriented Programming (which we will refer to as OOP) is a concept with numerous universally recognized virtues today. It's a programming method that helps improve application and software development and maintenance, with significant time savings. It's important to keep in mind that it doesn't reject structured (or procedural) programming since it's built upon it. The classic approach to introducing OOP into an existing language is through encapsulation of existing functionality. The C++ language was built on C and brought OOP. The JAVA language is based on C++ syntax.

Perl has also followed suit by offering extensions to give its fans the ability to use this programming paradigm. Nevertheless, Perl being permissive, it's not as strict as "pure object" languages. But this is normal, as Perl's philosophy is maintained, "there is more than one way to do it" (TIMTOWTDI).

I was strongly inspired by this tutorial[^1], but I'll only use a part of it for quick implementation. If you want a more condensed version that nevertheless requires good foundations, there is this one[^2].

## Simple Example

Let's start with a simple example including a module (because it's mandatory) and we'll see how to send information to it:

(`Personne.pm`)

```perl
# Nom du package, de notre classe
package Personne;
# Avertissement des messages d'erreurs
use warnings;
# Vérification des déclarations
use strict;

sub new {
  my ( $classe, $ref_arguments ) = @_;

  # Vérifions la classe
  $classe = ref($classe) || $classe;

  # Création de la référence anonyme d'un hachage vide (futur objet)
  my $this = {};

  # Liaison de l'objet à la classe
  bless( $this, $classe );

  $this->{_NOM}          = $ref_arguments->{nom};
  $this->{_PRENOM}       = $ref_arguments->{prenom};
  $this->{AGE}           = $ref_arguments->{age};
  $this->{_SEXE}         = $ref_arguments->{sexe};
  $this->{NOMBRE_ENFANT} = $ref_arguments->{nombre_enfant};

  return $this;
}

# Méthode marcher - ne prend aucun argument
sub marcher {
  my $this = shift;

  print "[$this->{_NOM} $this->{_PRENOM}] marche\n";

  return;
}

1;                # Important, à ne pas oublier
__END__           # Le compilateur ne lira pas les lignes après elle
```

And then you have your main program:

(`soft.pl`)

```perl
#!/usr/bin/perl
use warnings;
use strict;

use Personne;

my $Objet_Personne1 = Personne->new(
  { nom           => 'Dupont',
    prenom        => 'Jean',
    age           => 45,
    sexe          => 'M',
    nombre_enfant => 3,
  }
);
```

## References

[^1]: http://djibril.developpez.com/tutoriels/perl/poo/
[^2]: http://woufeil.developpez.com/tutoriels/perl/poo/
