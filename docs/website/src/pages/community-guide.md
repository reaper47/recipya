# Community Guide

Welcome to the project's community guide! 

## Support Channels

- [Matrix space](https://app.element.io/#/room/#recipya:matrix.org)
- [GitHub Issues](https://github.com/reaper47/recipya/issues)
- [GitHub Discussions](https://github.com/reaper47/recipya/discussions)

## Ways of Contributing

Recipya stands as a collaborative open-source endeavor. We welcome anyone who wants to help
us make this recipes manager the best it can be! Your input and contributions are vital as we work towards creating
an amazing recipe management solution. Not knowing how to code is not a requirement to contribute!

The following subsections elaborate on the different ways you can contribute.

### Feature Development

Please feel free to work on features that are [unassigned](https://github.com/reaper47/recipya/issues?q=is%3Aopen+is%3Aissue+label%3Aenhancement+no%3Aassignee).
I also encourage you to open a [feature request](https://github.com/reaper47/recipya/issues/new?assignees=&labels=enhancement&projects=&template=feature_request.md&title=)
issue when you have ideas that may improve the software.

You are not required to implement features yourself if you feel uncomfortable. However, the process is as follows if you 
do wish to.

1. Check the list of features currently [requested](https://github.com/reaper47/recipya/issues?q=is%3Aopen+is%3Aissue+label%3Aenhancement+no%3Aassignee).
2. Select the one you want to work on.
3. Comment that you want to fix it or send me a message in the [Matrix](https://app.element.io/#/room/#recipya:matrix.org)
   space so that I can move the task to the `in progress` column in the [board](https://github.com/users/reaper47/projects/2)
   and assign you to it.
4. Fork the repository if you have not done so.
5. Implement the feature and write tests.
6. Push the changes to your fork.
7. Open a pull request so that I can merge your work into `main`.

:::note 

Please be aware that working on a feature without opening an issue on GitHub first might lead to rejection if I believe
it is not a good fit for Recipya.

:::

### Bugs

Feel free to file bugs when you discover some! Please ensure that the bug you found has not been reported 
in the [GitHub issues](https://github.com/reaper47/recipya/issues?q=is%3Aopen+is%3Aissue+label%3Abug) before filing 
an [issue](https://github.com/reaper47/recipya/issues/new?assignees=&labels=&projects=&template=bug_report.md&title=) 
to reduce the number of duplicates.

You are not required to fix the bug yourself if you feel uncomfortable. However, the process is as follows if you 
do wish to.

1. Check the list of bugs currently [filed](https://github.com/reaper47/recipya/issues?q=is%3Aopen+is%3Aissue+label%3Abug).
2. Select an issue you want to work on.
3. Comment that you want to fix it or send me a message in the [Matrix](https://app.element.io/#/room/#recipya:matrix.org) 
   space so that I can move the task to the `in progress` column in the [board](https://github.com/users/reaper47/projects/2)
   and assign you to it.
4. Fork the repository if you have not done so.
5. Fix the bug, test it properly, and push the changes to your fork.
6. Open a pull request so that I can merge your work into `main`.

### Documentation

This website is the official documentation for Recipya. It is built using [Docusaurus](https://docusaurus.io/), which
is a static site generator that is very easy to use and understand. You do not need to open an issue regarding updates 
to the documentation. Please feel free to update as you deem fit and open a pull request. You can help us with 
translations, adding a language, fixing typos, improving grammar, adding sections, updating images, versioning, etc.

To develop the documentation locally, you must first [fork](https://github.com/reaper47/recipya/fork) the project.

Then, open a command prompt or terminal and navigate to `recipya/docs/website`.

```bash
cd path/to/recipya/docs/website
```

Next, run the following command to install the node modules. 
Please install [npm](https://nodejs.org/en/download) if the command is not found.

```bash
npm install
```

Finally, start the local website. 

```bash
npm run start
```

The website should have opened in your browser automatically at http://localhost:3000. 
You are now free to edit text and changes will be reflected in the browser on file save.

### Helping Others

It is always great to help anyone who needs a hand. Please see the [support channels](/community-guide#support-channels)
for places where you could lend a hand.
