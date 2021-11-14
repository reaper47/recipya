import { Fragment, useState } from 'react';
import { NavigateFunction, useNavigate } from 'react-router-dom';
import { Dropdown, DropdownItemProps, Icon } from 'semantic-ui-react';
import { AboutPortal } from './AboutPortal';

export const Avatar = (): JSX.Element => {
  const navigate: NavigateFunction = useNavigate();

  const [isAboutOpen, setIsAboutOpen] = useState<boolean>(false);
  const handleAboutClose = () => setIsAboutOpen(false);
  const handleAboutOpen = () => setIsAboutOpen(true);

  const to = (
    _event: React.MouseEvent<HTMLDivElement, MouseEvent>,
    data: DropdownItemProps
  ) => navigate(data.value as string);

  const renderMenu = (): JSX.Element => {
    const userIcon: JSX.Element = <Icon link name='user circle' size='big' />;

    return (
      <Dropdown className='icon' icon={null} trigger={userIcon} floating>
        <Dropdown.Menu>
          <Dropdown.Header icon='user' content='Adam' />
          <Dropdown.Divider />

          <Dropdown.Item
            icon='settings'
            content='Settings'
            value='/settings'
            onClick={to}
          />
          <Dropdown.Divider />

          <Dropdown.Item
            icon='info circle'
            content='About'
            onClick={handleAboutOpen}
          />

          <Dropdown.Item icon='power' content='Log out' />
        </Dropdown.Menu>
      </Dropdown>
    );
  };

  return (
    <Fragment>
      {renderMenu()}
      <AboutPortal isOpen={isAboutOpen} handleClose={handleAboutClose} />
    </Fragment>
  );
};
