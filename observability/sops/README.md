### SOP writing guidelines
When writing SOPS for different topics please try to keep the following guidelines in mind and use the following [template](https://github.com/bf2fc6cc711aee1a0c2a/cos-tools/blob/main/observability/sops/templates/alert.asciidoc). This is to ensure the engineers who will be using these SOPS in a production scenario (potentially from being woken up) have no issues.

#### 1. Concise instructions

When writing these SOPS please avoid writing long, wordy paragraphs. The main purpose of these SOPS is to inform the reader on how to solve the issue at hand in a quick and easy to follow manner.

When writing steps for *Execute/Resolution*, the preferred format is the following:

1. Add a user to cluster (quick explanation of what your about to do)
   `command blah blah` (command to do it)

#### 2. Terms and spelling

Please use the following terms when referring to resources (outside of code blocks) as well as keeping an eye on spelling.

* OSD Data Plane cluster
* Kafka (Capital K)

#### 3. Links

Please add links to resources where possible this allows easy access to information without too much searching.
For example Adding links to credentials in the prerequisite.
Format example:
* [Link](http://not-real-link-example)

**Rather than**

* http://not-real-link-example

#### 4. Code blocks

Where possible please keep code blocks on new lines rather than as part of a sentence. For example

1. To get the users for the following cluster:

   `command blah blah`

**Rather than**

1. To get users for the following cluster run `command blah blah`