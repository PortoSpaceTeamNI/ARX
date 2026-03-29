import type { IDockviewHeaderActionsProps } from 'dockview-react';
import { Minus, PictureInPicture2 } from 'lucide-react';
import { useEffect, useState } from 'react';

import Button from '@/elements/Button';

export default function PopoutButton({
  api,
  group,
  containerApi,
}: IDockviewHeaderActionsProps) {
  const [isPopout, setIsPopout] = useState(api.location.type === 'popout');

  useEffect(() => {
    const disposable = group.api.onDidLocationChange((event) => {
      setIsPopout(event.location.type === 'popout');
    });

    return () => disposable.dispose();
  }, [group.api]);

  function togglePopout() {
    if (isPopout) {
      const group = containerApi.addGroup();
      group.api.moveTo({ group });
    } else {
      containerApi.addPopoutGroup(group, {
        popoutUrl: '/popout.html',
      });
    }
  }

  return (
    <Button variant="ghost" size="icon" onClick={togglePopout}>
      {isPopout ? <Minus /> : <PictureInPicture2 />}
    </Button>
  );
}
