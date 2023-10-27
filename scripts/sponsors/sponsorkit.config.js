import {defineConfig} from 'sponsorkit';

const tiers = {
    helper: {
        avatar: {
            size: 45
        },
        boxWidth: 55,
        boxHeight: 55,
        container: {
            sidePadding: 30
        },
    },
    coffeeBoost: {
        avatar: {
            size: 50
        },
        boxWidth: 65,
        boxHeight: 65,
        container: {
            sidePadding: 30
        },
    },
    bronzeWarrior: {
        avatar: {
            size: 55
        },
        boxWidth: 75,
        boxHeight: 75,
        container: {
            sidePadding: 20
        },
        name: {
            maxLength: 10
        },
    },
    silverRainbow: {
        avatar: {
            size: 65
        },
        boxWidth: 90,
        boxHeight: 80,
        container: {
            sidePadding: 30
        },
        name: {
            maxLength: 10
        }
    },
    goldenBridge: {
        avatar: {
            size: 85
        },
        boxWidth: 110,
        boxHeight: 100,
        container: {
            sidePadding: 30
        },
        name: {
            maxLength: 20
        }
    },
}

const composeFunc = (composer, _tierSponsors, _config) => {
    composer.addSpan(20);
}

export default defineConfig({
    github: {
        login: 'reaper47',
        type: 'user',
    },

    width: 800,
    formats: ['svg'],
    tiers: [
        {
            title: 'Helpers',
            preset: tiers.helper,
            composeAfter: composeFunc,
        },
        {
            title: 'Coffee Boosters',
            monthlyDollars: 5,
            preset: tiers.coffeeBoost,
            composeAfter: composeFunc,
        },
        {
            title: 'Bronze Warriors',
            monthlyDollars: 15,
            preset: tiers.bronzeWarrior,
            composeAfter: composeFunc,
        },
        {
            title: 'Silver Rainbows',
            monthlyDollars: 20,
            preset: tiers.silverRainbow,
            composeAfter: composeFunc,
        },
        {
            title: 'Golden Bridges',
            monthlyDollars: 50,
            preset: tiers.goldenBridge,
            composeAfter: composeFunc,
        },
    ],
});