import {
  ActionIcon,
  Center,
  Popover,
  SimpleGrid,
} from '@mantine/core';
import { useDisclosure } from '@mantine/hooks';
import { ListIcon } from './ListIcon';
import { listIconCatalog, type ListIconId } from './icon-catalog';

type IconPickerProps = {
  onChange: (value: ListIconId) => void;
  value: ListIconId;
};

export function IconPicker({ onChange, value }: IconPickerProps) {
  const [opened, { close, open, toggle }] = useDisclosure(false);

  function selectIcon(iconId: ListIconId) {
    onChange(iconId);
    close();
  }

  return (
    <Popover
      opened={opened}
      onChange={(nextOpened) => (nextOpened ? open() : close())}
      position="bottom-start"
      width={260}
      shadow="md"
      trapFocus
      returnFocus
    >
      <Popover.Target>
        <ActionIcon
          variant="light"
          size="xl"
          radius="sm"
          aria-label="Choose list icon"
          onClick={toggle}
        >
          <ListIcon iconId={value} size={22} />
        </ActionIcon>
      </Popover.Target>

      <Popover.Dropdown p="xs">
        <SimpleGrid cols={6} spacing={4}>
          {listIconCatalog.map((option) => {
            const selected = option.id === value;

            return (
              <Center key={option.id}>
                <ActionIcon
                  variant={selected ? 'filled' : 'light'}
                  size={36}
                  radius="sm"
                  aria-label={`Use ${option.label} icon`}
                  aria-pressed={selected}
                  onClick={() => selectIcon(option.id)}
                >
                  <ListIcon iconId={option.id} size={22} />
                </ActionIcon>
              </Center>
            );
          })}
        </SimpleGrid>
      </Popover.Dropdown>
    </Popover>
  );
}
