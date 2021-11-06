import { Fragment, MouseEvent, useState } from 'react';
import { Icon } from 'semantic-ui-react';

export const Drawer = (): JSX.Element => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const handleClick = (event: MouseEvent<HTMLElement>): void => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = (): void => {
    setAnchorEl(null);
  };

  return (
    <Fragment>
      <Icon className='bars' fitted />
    </Fragment>
  );
};
