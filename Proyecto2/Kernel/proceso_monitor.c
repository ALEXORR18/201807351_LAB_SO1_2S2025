#include <linux/init.h>
#include <linux/module.h>
#include <linux/proc_fs.h>
#include <linux/uaccess.h>

#define PROC_NAME_SYSINFO "sysinfo_so1_201807351"
#define PROC_NAME_CONTINFO "continfo_so1_201807351"

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Kai");
MODULE_DESCRIPTION("Módulo kernel mínimo para Proyecto SO1");

static struct proc_dir_entry *sysinfo_entry;
static struct proc_dir_entry *continfo_entry;

#define BUF_SIZE 1024
static char sysinfo_buf[BUF_SIZE];
static char continfo_buf[BUF_SIZE];

static ssize_t sysinfo_read(struct file *file, char __user *user_buf, size_t count, loff_t *ppos)
{
    size_t len = snprintf(sysinfo_buf, BUF_SIZE, "{}\n"); // JSON vacío
    return simple_read_from_buffer(user_buf, count, ppos, sysinfo_buf, len);
}

static ssize_t continfo_read(struct file *file, char __user *user_buf, size_t count, loff_t *ppos)
{
    size_t len = snprintf(continfo_buf, BUF_SIZE, "[]\n"); // JSON vacío
    return simple_read_from_buffer(user_buf, count, ppos, continfo_buf, len);
}

static const struct proc_ops sysinfo_fops = {
    .proc_read = sysinfo_read,
};

static const struct proc_ops continfo_fops = {
    .proc_read = continfo_read,
};

static int __init proc_init(void)
{
    sysinfo_entry = proc_create(PROC_NAME_SYSINFO, 0444, NULL, &sysinfo_fops);
    if (!sysinfo_entry) {
        pr_err("Error creando /proc/%s\n", PROC_NAME_SYSINFO);
        return -ENOMEM;
    }

    continfo_entry = proc_create(PROC_NAME_CONTINFO, 0444, NULL, &continfo_fops);
    if (!continfo_entry) {
        pr_err("Error creando /proc/%s\n", PROC_NAME_CONTINFO);
        proc_remove(sysinfo_entry);
        return -ENOMEM;
    }

    pr_info("Módulo kernel cargado: archivos /proc creados\n");
    return 0;
}

static void __exit proc_exit(void)
{
    proc_remove(sysinfo_entry);
    proc_remove(continfo_entry);
    pr_info("Módulo kernel descargado: archivos /proc eliminados\n");
}

module_init(proc_init);
module_exit(proc_exit);

