import { useCombobox } from 'downshift';
import { useEffect, useRef, useState } from 'react';

interface ComboboxProps {
    type?: string,
    id?: string,
    name?: string,
    placeholder?: string,
    className?: string,
    value?: string,
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void,
    onClick?:  (e: React.MouseEvent<HTMLInputElement>) => void,
    onBlur?: () => void,
}

const sampleData = [
    "a",
    "aa",
    "aaa",
    "b",
    "c",
    "d"
];

export function Combobox(props: ComboboxProps) {
    const inputRef = useRef<HTMLInputElement | null>(null);
    const [offset, setOffset] = useState<{ left: number, width: number }>({ left: 0, width: 0 });
    const [items, setItems] = useState<string[]>([])
    const { isOpen,
        getToggleButtonProps,
        getLabelProps,
        getMenuProps,
        getInputProps,
        highlightedIndex,
        getItemProps,
        selectedItem
    } = useCombobox({
        onInputValueChange({ inputValue }) {
             setItems(sampleData.filter(i => !inputValue || i.includes(inputValue.toLowerCase())))
        },
        items,
    });

    useEffect(() => {
        function updateBox() {
            if (inputRef.current) {
                const rect = inputRef.current.getBoundingClientRect();
                console.log(rect);
                setOffset({
                    left: rect.left,
                    width: rect.width,
                });
            }
        }
        updateBox();
        window.addEventListener('resize', updateBox);
        return () => window.removeEventListener('resize', updateBox);
    }, []);

    return (
        <div className={`min-w-0 ${props.className} relative`}>
            <input
                type={props.type || "text"}
                name={props.name || ""}
                placeholder={props.placeholder || ""}
                className={`w-full border-0 border-b-2 border-blue-200 focus:border-blue-400 focus:outline-none focus:shadow-none`}
                onChange={props.onChange}
                onBlur={props.onBlur}
                {...getInputProps({
                    id: props.id || "",
                    onClick: props.onClick,
                    ref: node => {inputRef.current = node;}
                })}
            />
            <ul
                className={`fixed bg-black mt-1 max-h-80 overflow-scroll p-0 z-10 border-2 border-green-200 border-t-0 ${(isOpen && items.length) ? '' : 'hidden'}`}
                style={{ width: `${offset.width}px`, left:`${offset.left}px` }}
                {...getMenuProps()}
            >
                {isOpen && items.map((item, index) => (
                    <li
                        key={index}
                        className={`${highlightedIndex === index && 'bg-blue-300'} ${selectedItem === item && 'font-bold'} py-2 px-3 shadow-sm flex flex-col select-none hover:cursor-pointer hover:bg-gray-800`}
                        {...getItemProps({item, index})}
                    >
                        {item}
                    </li>
                ))}
            </ul>
        </div>
    );
}
