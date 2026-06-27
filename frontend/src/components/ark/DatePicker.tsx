import { DatePicker as ArkDatePicker, DateInputValueChangeDetails } from '@ark-ui/solid'
import { Index } from 'solid-js'
import { Portal } from 'solid-js/web'
import { parseDate } from '@ark-ui/solid';
import 'src/styles/datepicker.css';

interface DatePickerProps {
    onChange: (d: DateInputValueChangeDetails) => void
    value: Date
}

export default function DatePicker(props: DatePickerProps) {
    return (
        <ArkDatePicker.Root value={props.value ? [parseDate(props.value)] : []} onValueChange={props.onChange}>
            <ArkDatePicker.Context>
                {(api) => (
                    <>
                        <ArkDatePicker.Label>Date</ArkDatePicker.Label>
                        <ArkDatePicker.Control>
                            <ArkDatePicker.Input onFocus={() => api().setOpen(true)} />
                        </ArkDatePicker.Control>
                        <Portal>
                            <ArkDatePicker.Positioner>
                                <ArkDatePicker.Content>
                                    <ArkDatePicker.View view="day">
                                        <ArkDatePicker.Context>
                                            {(api) => (
                                                <>
                                                    <ArkDatePicker.ViewControl>
                                                        <ArkDatePicker.PrevTrigger>‹</ArkDatePicker.PrevTrigger>
                                                        <ArkDatePicker.MonthSelect />
                                                        <ArkDatePicker.YearSelect />
                                                        <ArkDatePicker.NextTrigger>›</ArkDatePicker.NextTrigger>
                                                    </ArkDatePicker.ViewControl>
                                                    <ArkDatePicker.Table>
                                                        <ArkDatePicker.TableHead>
                                                            <ArkDatePicker.TableRow>
                                                                <Index each={api().weekDays}>{(d) => <ArkDatePicker.TableHeader>{d().short}</ArkDatePicker.TableHeader>}</Index>
                                                            </ArkDatePicker.TableRow>
                                                        </ArkDatePicker.TableHead>
                                                        <ArkDatePicker.TableBody>
                                                            <Index each={api().weeks}>{(week) =>
                                                                <ArkDatePicker.TableRow>
                                                                    <Index each={week()}>{(day) =>
                                                                        <ArkDatePicker.TableCell value={day()}>
                                                                            <ArkDatePicker.TableCellTrigger>{day().day}</ArkDatePicker.TableCellTrigger>
                                                                        </ArkDatePicker.TableCell>
                                                                    }</Index>
                                                                </ArkDatePicker.TableRow>
                                                            }</Index>
                                                        </ArkDatePicker.TableBody>
                                                    </ArkDatePicker.Table>
                                                </>
                                            )}
                                        </ArkDatePicker.Context>
                                    </ArkDatePicker.View>
                                </ArkDatePicker.Content>
                            </ArkDatePicker.Positioner>
                        </Portal>
                    </>
                )}
            </ArkDatePicker.Context>
        </ArkDatePicker.Root>
    )
}
