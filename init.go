package gokube

func (k *Kube) Init() error {
	err := k.GetClientset()
	if err != nil {
		return err
	}
	return nil
}
