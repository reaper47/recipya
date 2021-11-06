import { Fragment, useState } from 'react';
import { NavigateFunction, useNavigate } from 'react-router-dom';
import { Dropdown, DropdownItemProps, Icon } from 'semantic-ui-react';
import { TablePortal } from '../portals/TablePortal';

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

  const renderAboutPortal = (): JSX.Element => {
    const rows = [
      {
        cells: [
          'Version:',
          <a
            className='green-a'
            href='https://github.com/reaper47/recipya/releases/tag/v1.0.0'
            target='_blank'
          >
            1.0.0
          </a>,
        ],
      },
      {
        cells: [
          'Home page:',
          <a
            className='green-a'
            href='https://www.musicavis.ca'
            target='_blank'
          >
            www.musicavis.ca
          </a>,
        ],
      },
      {
        cells: [
          'Source code:',
          <a
            className='green-a'
            href='https://github.com/reaper47/recipya'
            target='_blank'
          >
            github.com/reaper47/recipya
          </a>,
        ],
      },
      {
        cells: [
          'Feature requests:',
          <a
            className='green-a'
            href='https://github.com/reaper47/recipya/issues'
            target='_blank'
          >
            github.com/reaper47/recipya/issues
          </a>,
        ],
      },
    ];

    return (
      <TablePortal
        title='Recipya Recipes Manager'
        rows={rows}
        isOpen={isAboutOpen}
        handleClose={handleAboutClose}
      />
    );
  };

  return (
    <Fragment>
      {renderMenu()}
      {renderAboutPortal()}
    </Fragment>
  );
};
