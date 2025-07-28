interface InputProps {
    type?: string,
    id?: string,
    name?: string,
    placeholder?: string,
    className?: string,
    value?: string,
    onChange: (e: React.ChangeEvent<HTMLInputElement>) => void,
}

export function Input(props: InputProps) {
    return (
        <input
            type={props.type || "text"}
            id={props.id || ""}
            name={props.name || ""}
            placeholder={props.placeholder || ""}
            className={`min-w-0 border-0 border-b-2 border-blue-200 focus:border-blue-400 focus:outline-none focus:shadow-none ${props.className}`}
            value={props.value || ""}
            onChange={props.onChange}
        />
    );
}
