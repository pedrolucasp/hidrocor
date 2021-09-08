# hidrocor

A simple, tiny, markdown-based wiki engine. Useful for creating simple wikis,
with minimal effort on deployment, configuration and management.

## Preamble

hidrocor should be installed via your system package manager. You can, however,
compile from source by grabbing a tarball.

If installed from the source, you just need to run:

```sh
$ make
# make install
```

## Configurating wikis

A wiki, is just a folder containing some markdown files. hidrocor will parse and
turn them into simple pages, and link any `/*.md` contained inside the same
folder. You can mantain these wikis as a git repository for collaboration.

By default, hidrocor will lookup a configuration file in `/etc/hidrocor.yml`.
You can supply one specific configuration file with the `-c` parameter.

## Contributing

The upstream repository can be found [here][repo]. Send patches to
[porcellis@eletrotupi.com][mailing-list] or open pull requests [on the GitHub
mirror][github].

## License

GNU GPL-3.0

[repo]: https://git.eletrotupi.com/hidrocor
[mailing-list]: mailto:porcellis@eletrotupi.com
[github]: https://github.com/pedrolucasp/hidrocor
