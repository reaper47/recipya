import { TablePortal } from '../portals/TablePortal';

interface AboutPortalProps {
  isOpen: boolean;
  handleClose(): void;
}

export const AboutPortal = ({
  isOpen,
  handleClose,
}: AboutPortalProps): JSX.Element => {
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
        <a className='green-a' href='https://www.musicavis.ca' target='_blank'>
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
      isOpen={isOpen}
      handleClose={handleClose}
    />
  );
};
