import {defineConfig, presets} from 'sponsorkit';

export default defineConfig({
    github: {
        login: 'reaper47',
        type: 'user',
    },

    width: 800,
    formats: ['svg'],
    tiers: [
        {
            title: 'Past Sponsors',
            monthlyDollars: -1,
            preset: presets.xs,
        },
        // Default tier
        {
            title: 'Backers',
            preset: presets.base,
        },
        {
            title: 'Sponsors',
            monthlyDollars: 10,
            preset: presets.medium,
        },
        {
            title: 'Silver Sponsors',
            monthlyDollars: 50,
            preset: presets.large,
        },
        {
            title: 'Gold Sponsors',
            monthlyDollars: 100,
            preset: presets.xl,
        },
    ],
});