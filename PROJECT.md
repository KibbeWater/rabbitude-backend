# Project Plans

The project is grouped into three categories: High Priority, Medium Priority, and Low Priority.

What are seen as "High Priority" are features that are required before an initial "1.0" release.

"Medium Priority" features are what could be seen as important to a large number of users, but there exists alternatives or workarounds.

"Low Priority" features are what are seen as nice-to-haves, but are not required for the initial release.

If you have any suggestions or ideas, you are free to modify this document and submit a pull request.

## High Priority

### General

- [ ] Add Tests and CI/CD
- [ ] Add Documentation
- [ ] Overhauled Config system, larger configuration alternatives and more functional setup screens
- [ ] Make select custom services opt-out without having to deny on restart

### Providers

- [ ] Add Groq Provider (LLMs / Speech services)
- [ ] Add OpenAI Providers (GPT-3.5/4/4o, Whisper, etc.)
- [ ] Add a "Search" Provider (Perplexity?)
- [ ] On-Device Text-to-Speech Provider

## Medium Priority

- [ ] Plugin System, "deeper" integration with the program, but does not act as a "provider"
- [ ] Add more commercial providers (Google, Amazon, etc.)
- [ ] Add a "Vision" service provider (GPT-4o)

## Low Priority

- [ ] Add more of Rabbits core integrations (Music, Rideshare, Food ordering, etc.)
- [ ] Make the Setup screens prettier, colors, emojis!
- [ ] **[FUTURE]** Apple Intelligence Provider, support for Apple's Greymatter APIs in September

## Features

### Provider(s)

Providers are essentially different "services", they are used quite interchangeably throughout the project. Providers provide different services such as "LLM, "Speech", "Search", etc.

Prompts come throgh and gets tagged with an "Intent", which defines which service to use. Providers can then expose different services (custom or included as "base" services) to be used.

**Now, what is the difference between a custom service and a base service?**

It might not seem very intuitive, but a base service are "essential" services that provide a baseline experience as expected by the user. They are more well-defined in the project to give the system a larger control over how each and every of those services work.

Custom services are services that are not included in the base experience, but can be added by one or more provider and will automatically appear setup screens.

## Setup / Configuration

The setup screens are designed to be as user-friendly as possible. However are a bit glitchy, below are a few things that could see improvements:

- [ ] Save / Restore console logs when entering and exiting setup screens
- [ ] Setup will require some sort of pagination, as the current system will not be able to handle more than 9 providers on one service
- [ ] Perhaps add more of a "GUI" style to the terminal, allowing for navigation, more advanced statistics and information whilst running the program
- [ ] General improvements to *looks* to the setup screens with more colors, emojis. Just a bit more life
