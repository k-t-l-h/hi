package repository

const (
	CREATEPERSONA  = "INSERT INTO public.persona(name, age, address, workplace) VALUES ($1, $2, $3, $4) RETURNING id;"
	READPERSONA    = "SELECT name, age, address, workplace FROM public.persona WHERE id = $1;"
	UPDATEPERSONA  = "UPDATE public.persona SET name=$1, age=$2, address=$3, workplace=$4 WHERE id = $5;"
	DELETEPERSONA  = "DELETE FROM public.persona WHERE id = $1;"
	READALLPERSONA = "SELECT id, name, age, address, workplace FROM public.persona;"
)
