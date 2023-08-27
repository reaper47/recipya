import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';

const FeatureList = [
  {
    title: 'Curate recipes',
    Svg: require('@site/static/img/books.svg').default,
    description: (
        <>
            Effortlessly import your favorite recipes from around the web,
            digitize paper recipes, and add recipes manually.
        </>
    ),
  },
  {
    title: 'Stick to one measurement system',
    Svg: require('@site/static/img/metric-seal.svg').default,
    description: (
      <>
        All of your recipes can be converted to your preferred measurement system.
        Say goodbye to imperial if you are a metric person.
      </>
    ),
  },
  {
    title: 'Self-hostable',
    Svg: require('@site/static/img/devices.svg').default,
    description: (
      <>
        Easily self-host the software on your server with Docker.
      </>
    ),
  },
];

function Feature({Svg, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div>
      <div className="text--center padding-horiz--md">
        <h3>{title}</h3>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
