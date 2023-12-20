---
title: Wakelock
weight: 8
---

Imagine you are cooking a recipe from a mobile device. You gather all the ingredients, and you are ready to start, 
but your device goes to sleep. You unlock your device and start getting your hands dirty in the flour. Your device 
goes to sleep again. You clean your hands, unlock your phone, and continue the recipe. Your device goes to sleep 
once again. You are frustrated. This process continues over and over. Your device is dirty and so is your food.

You want the device not to sleep while you are cooking.

Recipya solves this problem by providing a browser wakelock when viewing a recipe, preventing the device from 
going to sleep.

:::note

The wakelock is turned on automatically when viewing a recipe.

:::

On supported browsers, you will notice a light bulb icon to the far left of a recipe's title. When the light bulb 
is on, the wakelock is enabled and the screen will not sleep.

![img alt](/img/features/wakelock-on.png)

When turned off, the wakelock is disabled and the screen will eventually sleep.

![img alt](/img/features/wakelock-off.png)

:::caution

The wakelock feature is not supported on the following platforms:
- Firefox
- Firefox for Android

:::