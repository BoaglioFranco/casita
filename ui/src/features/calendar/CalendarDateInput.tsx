import { DatePickerInput, type DatePickerInputProps } from '@mantine/dates';
import { useMediaQuery } from '@mantine/hooks';
import { CalendarDays } from 'lucide-react';

type CalendarDateInputProps = DatePickerInputProps & {
  onPickerOpenChange?: (opened: boolean) => void;
};

export function CalendarDateInput({
  onPickerOpenChange,
  onClick,
  onChange,
  onDropdownClose,
  modalProps,
  ...props
}: CalendarDateInputProps) {
  const mobile = useMediaQuery('(max-width: 47.99em)');

  return (
    <DatePickerInput
      size="sm"
      radius="md"
      valueFormat="ddd, D MMM YYYY"
      firstDayOfWeek={1}
      weekendDays={[]}
      highlightToday
      styles={{ input: { minHeight: 44 } }}
      dropdownType={mobile ? 'modal' : 'popover'}
      onClick={(event) => {
        onPickerOpenChange?.(true);
        onClick?.(event);
      }}
      onChange={(value) => {
        onPickerOpenChange?.(false);
        onChange?.(value);
      }}
      onDropdownClose={() => {
        onPickerOpenChange?.(false);
        onDropdownClose?.();
      }}
      rightSection={<CalendarDays size={18} aria-hidden />}
      rightSectionPointerEvents="none"
      popoverProps={{
        withinPortal: false,
        shadow: 'md',
        radius: 'md',
        transitionProps: { transition: 'pop', duration: 150 },
      }}
      modalProps={{
        title: 'Choose date',
        centered: true,
        size: 'auto',
        radius: 'md',
        padding: 'sm',
        withCloseButton: true,
        closeButtonProps: { 'aria-label': 'Close date picker' },
        transitionProps: { transition: 'pop', duration: 180 },
        ...modalProps,
        onExitTransitionEnd: () => {
          onPickerOpenChange?.(false);
          modalProps?.onExitTransitionEnd?.();
        },
      }}
      {...props}
    />
  );
}
