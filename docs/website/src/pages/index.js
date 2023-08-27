import React from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';

import styles from './index.module.css';
const Carousel = require("react-responsive-carousel").Carousel;
import 'react-responsive-carousel/lib/styles/carousel.min.css';

function HomepageHeader() {
    const {siteConfig} = useDocusaurusContext();
    return (
        <header className={clsx('hero hero--primary', styles.heroBanner)}>
            <div className="container">
                <h1 className="hero__title">{siteConfig.title}</h1>
                <p className="hero__subtitle">{siteConfig.tagline}</p>
                <div className={styles.buttons}>
                    <Link
                        className="button button--secondary button--lg"
                        to="/docs/intro">
                        Learn More
                    </Link>
                    <Link
                        className="button button--primary button--lg"
                        to="/docs/category/installation">
                        Get Started
                    </Link>
                </div>
            </div>
            <Carousel
                showArrows={false}
                width={"400px"}
                showThumbs={false}
                stopOnHover={false}
                showStatus={false}
                autoPlay={true}
                showIndicators={false}
                infiniteLoop={true}
                interval={4000}
                transitionTime={500}
                centerMode={true}
            >
                <div className="slide-item-box">
                    <img src="/img/carousel/collection.png" alt="All recipes"/>
                </div>
                <div className="slide-item-box">
                    <img src="/img/carousel/view.png" alt="Add recipe"/>
                </div>
                <div className="slide-item-box">
                    <img src="/img/carousel/add-recipe.png" alt="View recipe"/>
                </div>
            </Carousel>
        </header>
    );
}

export default function Home() {
    return (
        <Layout
            title="Welcome"
            description="A recipe manager for unforgettable family recipes, empowering you to curate and share your favorite recipes.">
            <HomepageHeader/>
            <main>
                <HomepageFeatures/>
            </main>
        </Layout>
    );
}
